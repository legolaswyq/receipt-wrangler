package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type capturedCustomOcrRequest struct {
	Method      string
	ContentType string
	Body        map[string]interface{}
}

func newMockCustomOcrServer(t *testing.T, statusCode int, responseBody string) (*httptest.Server, *capturedCustomOcrRequest) {
	t.Helper()
	captured := &capturedCustomOcrRequest{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.Method = r.Method
		captured.ContentType = r.Header.Get("Content-Type")

		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		captured.Body = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	t.Cleanup(server.Close)
	return server, captured
}

func customOcrSuccessBody(content string) string {
	resp := structs.OllamaTextResponse{}
	resp.Message.Role = "assistant"
	resp.Message.Content = content
	resp.Done = true
	raw, _ := json.Marshal(resp)
	return string(raw)
}

func newCustomOcrService(url, model string) OcrService {
	customOcr := models.CUSTOM
	return OcrService{
		ReceiptProcessingSettings: models.ReceiptProcessingSettings{
			OcrEngine:      &customOcr,
			OcrEngineUrl:   url,
			OcrEngineModel: model,
		},
	}
}

func TestOcrService_ReadImageWithCustomOcr_HappyPath(t *testing.T) {
	server, captured := newMockCustomOcrServer(t, http.StatusOK, customOcrSuccessBody("Store Name\nItem 1  $1.00"))

	service := newCustomOcrService(server.URL, "glm-ocr:latest")
	text, err := service.ReadImageWithCustomOcr([]byte("fake-image-bytes"))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if text != "Store Name\nItem 1  $1.00" {
		utils.PrintTestError(t, text, "Store Name\nItem 1  $1.00")
	}

	if captured.Method != http.MethodPost {
		utils.PrintTestError(t, captured.Method, http.MethodPost)
	}
	if captured.ContentType != "application/json" {
		utils.PrintTestError(t, captured.ContentType, "application/json")
	}
	if captured.Body["model"] != "glm-ocr:latest" {
		utils.PrintTestError(t, captured.Body["model"], "glm-ocr:latest")
	}
	if captured.Body["stream"] != false {
		utils.PrintTestError(t, captured.Body["stream"], false)
	}

	messages, ok := captured.Body["messages"].([]interface{})
	if !ok || len(messages) != 1 {
		utils.PrintTestError(t, captured.Body["messages"], "one message")
		return
	}
	message := messages[0].(map[string]interface{})
	images, ok := message["images"].([]interface{})
	if !ok || len(images) != 1 {
		utils.PrintTestError(t, message["images"], "one base64 image")
	}
}

func TestOcrService_ReadImageWithCustomOcr_ErrorField(t *testing.T) {
	resp := structs.OllamaTextResponse{Error: "model not found"}
	raw, _ := json.Marshal(resp)
	server, _ := newMockCustomOcrServer(t, http.StatusOK, string(raw))

	service := newCustomOcrService(server.URL, "glm-ocr:latest")
	_, err := service.ReadImageWithCustomOcr([]byte("fake-image-bytes"))
	if err == nil {
		utils.PrintTestError(t, err, "expected error for model-not-found response")
	}
}

func TestOcrService_ReadImageWithCustomOcr_MalformedResponseBody(t *testing.T) {
	server, _ := newMockCustomOcrServer(t, http.StatusOK, "not json at all")

	service := newCustomOcrService(server.URL, "glm-ocr:latest")
	_, err := service.ReadImageWithCustomOcr([]byte("fake-image-bytes"))
	if err == nil {
		utils.PrintTestError(t, err, "expected unmarshal error")
	}
}

func TestOcrService_ReadImageWithCustomOcr_ConnectionError(t *testing.T) {
	server, _ := newMockCustomOcrServer(t, http.StatusOK, customOcrSuccessBody(""))
	url := server.URL
	server.Close()

	service := newCustomOcrService(url, "glm-ocr:latest")
	_, err := service.ReadImageWithCustomOcr([]byte("fake-image-bytes"))
	if err == nil {
		utils.PrintTestError(t, err, "expected connection error")
	}
}

func writeTestImage(t *testing.T, path string, width, height uint, format string) {
	t.Helper()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("white")

	if err := mw.NewImage(width, height, pw); err != nil {
		t.Fatalf("ImageMagick NewImage: %v", err)
	}
	if err := mw.SetImageFormat(format); err != nil {
		t.Fatalf("ImageMagick SetImageFormat %s: %v", format, err)
	}
	if err := mw.WriteImage(path); err != nil {
		t.Fatalf("ImageMagick WriteImage: %v", err)
	}
}

func dimensionsOf(t *testing.T, blob []byte) (uint, uint, string) {
	t.Helper()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageBlob(blob); err != nil {
		t.Fatalf("ImageMagick ReadImageBlob: %v", err)
	}
	return mw.GetImageWidth(), mw.GetImageHeight(), mw.GetImageFormat()
}

func TestOcrService_ReadRawImage_DownscalesOversizedImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oversized.png")
	writeTestImage(t, path, 4032, 3024, "png")

	service := OcrService{}
	blob, err := service.readRawImage(path)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	width, height, format := dimensionsOf(t, blob)
	if width > customOcrMaxDimension || height > customOcrMaxDimension {
		utils.PrintTestError(t, nil, "expected longest side to be capped at customOcrMaxDimension")
	}
	if width != customOcrMaxDimension {
		utils.PrintTestError(t, width, customOcrMaxDimension)
	}
	// 4032x3024 is a 4:3 image, so height should scale proportionally.
	expectedHeight := uint(float64(3024) * (float64(customOcrMaxDimension) / float64(4032)))
	if height != expectedHeight {
		utils.PrintTestError(t, height, expectedHeight)
	}
	if format != "JPEG" {
		utils.PrintTestError(t, format, "JPEG")
	}
}

func TestOcrService_ReadRawImage_ConvertsFormatWithoutUpscalingSmallImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "small.png")
	writeTestImage(t, path, 200, 260, "png")

	service := OcrService{}
	blob, err := service.readRawImage(path)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	width, height, format := dimensionsOf(t, blob)
	if width != 200 || height != 260 {
		utils.PrintTestError(t, []uint{width, height}, []uint{200, 260})
	}
	if format != "JPEG" {
		utils.PrintTestError(t, format, "JPEG")
	}
}

func TestOcrService_ReadRawImage_MissingFile(t *testing.T) {
	service := OcrService{}
	_, err := service.readRawImage(filepath.Join(os.TempDir(), "does-not-exist.jpg"))
	if err == nil {
		utils.PrintTestError(t, err, "expected error for missing file")
	}
}

func TestOcrService_ReadImageWithCustomOcr_InvalidUrl(t *testing.T) {
	service := newCustomOcrService("http://%zz-not-valid", "glm-ocr:latest")
	_, err := service.ReadImageWithCustomOcr([]byte("fake-image-bytes"))
	if err == nil {
		utils.PrintTestError(t, err, "expected request construction error")
	}
}
