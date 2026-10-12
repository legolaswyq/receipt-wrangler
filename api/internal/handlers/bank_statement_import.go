package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// PreviewBankStatement parses an uploaded bank statement PDF and returns its
// transactions for review. Nothing is created.
func PreviewBankStatement(w http.ResponseWriter, r *http.Request) {
	errMessage := "Error reading bank statement"
	token := structs.GetClaims(r)

	// The command is loaded before the handler is built because the group it
	// names is what the permission gate is resolved against.
	command := commands.BankStatementPreviewCommand{}
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errMessage, http.StatusBadRequest)
		return
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errMessage,
		Writer:           w,
		Request:          r,
		GroupId:          utils.UintToString(command.GroupId),
		GroupPermissions: []string{permissions.GroupReceiptsCreate},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			importService := services.NewBankStatementImportService(nil)
			preview, err := importService.Preview(token.UserId, command)
			if errors.Is(err, services.ErrUnsupportedBankStatement) {
				utils.WriteCustomErrorResponse(w, err.Error(), http.StatusBadRequest)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := json.Marshal(preview)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

// ImportBankStatement creates a receipt for each reviewed statement transaction.
func ImportBankStatement(w http.ResponseWriter, r *http.Request) {
	errMessage := "Error importing bank statement"
	token := structs.GetClaims(r)

	command := commands.BankStatementImportCommand{}
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errMessage, http.StatusBadRequest)
		return
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errMessage,
		Writer:           w,
		Request:          r,
		GroupId:          utils.UintToString(command.GroupId),
		GroupPermissions: []string{permissions.GroupReceiptsCreate},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			importService := services.NewBankStatementImportService(nil)
			result, err := importService.Import(token.UserId, command)
			if errors.Is(err, services.ErrBankStatementForbidden) {
				utils.WriteCustomErrorResponse(w, err.Error(), http.StatusForbidden)
				return 0, nil
			}
			if errors.Is(err, services.ErrBankStatementInvalidSelection) {
				utils.WriteCustomErrorResponse(w, err.Error(), http.StatusBadRequest)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := json.Marshal(result)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}
