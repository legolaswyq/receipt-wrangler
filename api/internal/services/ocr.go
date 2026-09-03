package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"time"
)

const customOcrPrompt = "Transcribe all text from this receipt image exactly as it appears, preserving line order and layout. Return only the raw transcribed text, no commentary or formatting."

type OcrService struct {
	BaseService
	ReceiptProcessingSettings models.ReceiptProcessingSettings
}

func NewOcrService(tx *gorm.DB, receiptProcessingSettings models.ReceiptProcessingSettings) OcrService {
	service := OcrService{
		BaseService: BaseService{
			DB: repositories.GetDB(),
			TX: tx,
		},
		ReceiptProcessingSettings: receiptProcessingSettings,
	}

	return service
}

func (service OcrService) ReadImage(path string) (string, commands.UpsertSystemTaskCommand, error) {
	var text string
	startTime := time.Now()
	systemTaskCommand := commands.UpsertSystemTaskCommand{
		Type:                 models.OCR_PROCESSING,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.RECEIPT_PROCESSING_SETTINGS,
		AssociatedEntityId:   service.ReceiptProcessingSettings.ID,
		StartedAt:            time.Now(),
	}

	isCustomOcr := service.ReceiptProcessingSettings.OcrEngine != nil && *service.ReceiptProcessingSettings.OcrEngine == models.CUSTOM

	var imageBytes []byte
	var err error
	if isCustomOcr {
		// A custom vision-model OCR engine, not a classical OCR engine — it wants the
		// natural image, not the bilevel/blurred/deskewed bytes prepareImage
		// produces for Tesseract/EasyOCR (which would degrade its accuracy).
		imageBytes, err = service.readRawImage(path)
	} else {
		imageBytes, err = service.prepareImage(path)
	}
	if err != nil {
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		return "", systemTaskCommand, err
	}

	if service.ReceiptProcessingSettings.OcrEngine != nil && *service.ReceiptProcessingSettings.OcrEngine == models.TESSERACT_NEW {
		text, err = service.ReadImageWithTesseract(imageBytes)
		if err != nil {
			systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
			systemTaskCommand.ResultDescription = err.Error()
			return "", systemTaskCommand, err
		}
	}

	if service.ReceiptProcessingSettings.OcrEngine != nil && *service.ReceiptProcessingSettings.OcrEngine == models.EASY_OCR_NEW {
		text, err = service.ReadImageWithEasyOcr(imageBytes)
		if err != nil {
			systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
			systemTaskCommand.ResultDescription = err.Error()
			return "", systemTaskCommand, err
		}
	}

	if isCustomOcr {
		text, err = service.ReadImageWithCustomOcr(imageBytes)
		if err != nil {
			systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
			systemTaskCommand.ResultDescription = err.Error()
			return "", systemTaskCommand, err
		}
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	logging.LogStd(logging.LOG_LEVEL_INFO, "OCR and Image processing took: ", elapsedTime)

	systemSettingsRepository := repositories.NewSystemSettingsRepository(service.TX)
	systemSettings, err := systemSettingsRepository.GetSystemSettings()
	if err != nil {
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		return "", systemTaskCommand, err
	}

	if systemSettings.DebugOcr {
		err = service.writeDebuggingFiles(text, path, imageBytes, elapsedTime)
		if err != nil {
			return "", commands.UpsertSystemTaskCommand{}, err
		}
	}

	ocrEndTime := time.Now()
	systemTaskCommand.EndedAt = &ocrEndTime
	systemTaskCommand.ResultDescription = text

	return text, systemTaskCommand, nil
}

func (service OcrService) ReadImageWithTesseract(preparedImageBytes []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetVariable("tessedit_char_blacklist", "!@#$%^&*()_+=-[]}{;:'\"\\|~`<>/?")
	if err != nil {
		return "", nil
	}

	err = client.SetImageFromBytes(preparedImageBytes)
	if err != nil {
		return "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}

func (service OcrService) ReadImageWithEasyOcr(preparedImageBytes []byte) (string, error) {
	fileRepository := repositories.NewFileRepository(nil)
	tempPath, err := fileRepository.WriteTempFile(preparedImageBytes)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempPath)

	var textBuffer bytes.Buffer
	var text string
	cmd := exec.Command("easyocr", "-l", "en", "-f", tempPath, "--detail", "0", "--gpu", "0", "--verbose", "0")
	cmd.Stdout = &textBuffer
	cmd.Stderr = io.Discard

	err = cmd.Run()
	if err != nil {
		return "", err
	}
	text = textBuffer.String()

	return text, nil
}

// customOcrMaxDimension caps the longest side sent to a custom vision-model OCR
// engine. Vision models spend a large, fixed chunk of their context window on
// image tokens regardless of content; an uncapped phone photo (e.g. 4032x3024)
// can consume nearly all of a model's default context, leaving almost no
// budget for the transcribed text and truncating the response.
const customOcrMaxDimension = 1600

func (service OcrService) readRawImage(path string) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(path)
	if err != nil {
		return nil, err
	}

	// Applies any EXIF orientation (common on phone photos, including HEIC)
	// so the image isn't sent to the model sideways or upside down.
	err = mw.AutoOrientImage()
	if err != nil {
		return nil, err
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if width > customOcrMaxDimension || height > customOcrMaxDimension {
		if width >= height {
			height = uint(float64(height) * (float64(customOcrMaxDimension) / float64(width)))
			width = customOcrMaxDimension
		} else {
			width = uint(float64(width) * (float64(customOcrMaxDimension) / float64(height)))
			height = customOcrMaxDimension
		}

		err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS)
		if err != nil {
			return nil, err
		}
	}

	// Vision models are generally trained/served on JPEG/PNG; converting away
	// from formats like HEIC avoids relying on the model's own image decoder.
	err = mw.SetImageFormat("jpeg")
	if err != nil {
		return nil, err
	}

	err = mw.SetImageCompressionQuality(85)
	if err != nil {
		return nil, err
	}

	return mw.GetImageBlob()
}

func (service OcrService) ReadImageWithCustomOcr(imageBytes []byte) (string, error) {
	body := map[string]interface{}{
		"model": service.ReceiptProcessingSettings.OcrEngineModel,
		"messages": []structs.AiClientMessage{
			{
				Role:    "user",
				Content: customOcrPrompt,
				Images:  []string{utils.Base64Encode(imageBytes)},
			},
		},
		"stream": false,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, service.ReceiptProcessingSettings.OcrEngineUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", constants.ApplicationJson)
	request.Close = true

	httpClient := http.Client{Timeout: constants.AiHttpTimeout}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseObject structs.OllamaTextResponse
	err = json.Unmarshal(responseBody, &responseObject)
	if err != nil {
		return "", err
	}

	if len(responseObject.Error) > 0 {
		return "", fmt.Errorf("custom ocr engine error: %s", responseObject.Error)
	}

	return responseObject.Message.Content, nil
}

func (service OcrService) writeDebuggingFiles(ocrText string, path string, imageBytes []byte, ocrDuration time.Duration) error {
	fileRepository := repositories.NewFileRepository(nil)
	pathParts := strings.Split(path, "/")
	filename := pathParts[len(pathParts)-1]

	tempPath := fileRepository.GetTempDirectoryPath()
	textFilePath := filepath.Join(tempPath, filename+".txt")
	imageFilePath := filepath.Join(tempPath, filename+".jpg")

	textBytes := []byte(ocrText)

	os.Remove(textFilePath)
	err := utils.WriteFile(textFilePath, textBytes)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return err
	}

	os.Remove(imageFilePath)
	imgFile, err := os.Create(imageFilePath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	err = jpeg.Encode(imgFile, img, nil)
	if err != nil {
		return err
	}

	err = utils.WriteFile(imageFilePath, imageBytes)
	if err != nil {
		return err
	}

	fmt.Println("OCR Text saved to: ", textFilePath)
	fmt.Println("OCR Image saved to: ", imageFilePath)
	fmt.Println("OCR and image processing duration: ", ocrDuration)

	return nil
}

func (service OcrService) prepareImage(path string) ([]byte, error) {
	mw := imagick.NewMagickWand()
	err := mw.ReadImage(path)
	if err != nil {
		return nil, err
	}

	err = mw.TrimImage(0)
	if err != nil {
		return nil, err
	}

	err = mw.SetImageType(imagick.IMAGE_TYPE_BILEVEL)
	if err != nil {
		return nil, err
	}

	err = mw.BlurImage(0, 1.5)
	if err != nil {
		return nil, err
	}

	err = mw.SharpenImage(0, 1)
	if err != nil {
		return nil, err
	}

	err = mw.EnhanceImage()
	if err != nil {
		return nil, err
	}

	err = mw.ContrastImage(false)
	if err != nil {
		return nil, err
	}

	err = mw.DeskewImage(.40)
	if err != nil {
		return nil, err
	}

	if service.ReceiptProcessingSettings.OcrEngine != nil &&
		*service.ReceiptProcessingSettings.OcrEngine == models.EASY_OCR_NEW {
		err = mw.ScaleImage(mw.GetImageWidth()/2, mw.GetImageHeight()/2)
		if err != nil {
			return nil, err
		}
	}

	return mw.GetImageBlob()
}
