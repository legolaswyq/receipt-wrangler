package wranglerasynq

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
)

type QuickScanTaskPayload struct {
	Token        *structs.Claims
	PaidByUserId uint
	GroupId      uint
	Status       models.ReceiptStatus
	CategoryIds  []uint
	TagIds       []uint
	Comment      string
	// TempPath / OriginalFileName carry a single image. TempPaths / OriginalFileNames carry several
	// images that form one long receipt (combine mode); when set they take precedence and produce a
	// single receipt with all images attached.
	TempPath          string
	OriginalFileName  string
	TempPaths         []string
	OriginalFileNames []string
}

func HandleQuickScanTask(context context.Context, task *asynq.Task) error {
	taskId, err := GetTaskIdFromContext(context)
	if err != nil {
		return HandleError(err)
	}

	var payload QuickScanTaskPayload

	err = json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return HandleError(err)
	}

	receiptService := services.NewReceiptService(nil)
	_, err = receiptService.QuickScan(services.QuickScanParams{
		Token:             payload.Token,
		PaidByUserId:      payload.PaidByUserId,
		GroupId:           payload.GroupId,
		Status:            payload.Status,
		CategoryIds:       payload.CategoryIds,
		TagIds:            payload.TagIds,
		Comment:           payload.Comment,
		TempPath:          payload.TempPath,
		OriginalFileName:  payload.OriginalFileName,
		TempPaths:         payload.TempPaths,
		OriginalFileNames: payload.OriginalFileNames,
		AsynqTaskId:       taskId,
	})
	if err != nil {
		return HandleError(err)
	}

	return nil
}
