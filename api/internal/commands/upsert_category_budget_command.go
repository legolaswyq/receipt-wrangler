package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

type UpsertCategoryBudgetCommand struct {
	CategoryId uint            `json:"categoryId"`
	Amount     decimal.Decimal `json:"amount"`
}

func (command *UpsertCategoryBudgetCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &command)
}

func (command *UpsertCategoryBudgetCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{}
	errorMap := make(map[string]string)

	if command.CategoryId == 0 {
		errorMap["categoryId"] = "Category is required"
	}
	if command.Amount.LessThanOrEqual(decimal.Zero) {
		errorMap["amount"] = "Amount must be greater than zero"
	}

	vErr.Errors = errorMap
	return vErr
}
