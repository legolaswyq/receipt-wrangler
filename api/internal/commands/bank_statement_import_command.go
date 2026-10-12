package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// MaxBankStatementTransactions caps one import so a crafted request cannot
// create an unbounded number of receipts in a single transaction.
const MaxBankStatementTransactions = 2000

// BankStatementPreviewCommand is the multipart upload of a statement PDF to be
// parsed for review. Nothing is created from it.
type BankStatementPreviewCommand struct {
	File    []byte
	GroupId uint
}

func (command *BankStatementPreviewCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(constants.MultipartFormMaxSize)
	if err != nil {
		return err
	}

	groupId, err := utils.StringToUint(strings.TrimSpace(r.FormValue("groupId")))
	if err != nil {
		return err
	}
	command.GroupId = groupId

	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	command.File, err = io.ReadAll(file)
	return err
}

func (command *BankStatementPreviewCommand) Validate() structs.ValidatorError {
	errors := make(map[string]string)

	if command.GroupId == 0 {
		errors["groupId"] = "Group Id is required"
	}

	if len(command.File) == 0 {
		errors["file"] = "File is required"
	}

	return structs.ValidatorError{Errors: errors}
}

// BankStatementImportTransaction is one reviewed statement row to create a
// receipt from.
type BankStatementImportTransaction struct {
	Date       string          `json:"date"`
	Name       string          `json:"name"`
	Amount     decimal.Decimal `json:"amount"`
	IsIncome   bool            `json:"isIncome"`
	CategoryId *uint           `json:"categoryId"`
}

// ParsedDate returns the row's calendar date at UTC midnight, matching how
// receipt dates are stored.
func (transaction BankStatementImportTransaction) ParsedDate() (time.Time, error) {
	return time.Parse(structs.BankStatementDateFormat, transaction.Date)
}

type BankStatementImportCommand struct {
	GroupId      uint                             `json:"groupId"`
	PaidByUserId uint                             `json:"paidByUserId"`
	Transactions []BankStatementImportTransaction `json:"transactions"`
}

func (command *BankStatementImportCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &command)
}

func (command *BankStatementImportCommand) Validate() structs.ValidatorError {
	errors := make(map[string]string)

	if command.GroupId == 0 {
		errors["groupId"] = "Group Id is required"
	}

	if command.PaidByUserId == 0 {
		errors["paidByUserId"] = "Paid By User Id is required"
	}

	if len(command.Transactions) == 0 {
		errors["transactions"] = "At least one transaction is required"
	}

	if len(command.Transactions) > MaxBankStatementTransactions {
		errors["transactions"] = fmt.Sprintf("At most %d transactions can be imported at once", MaxBankStatementTransactions)
	}

	for i, transaction := range command.Transactions {
		basePath := fmt.Sprintf("transactions.%d", i)

		if len(strings.TrimSpace(transaction.Name)) == 0 {
			errors[basePath+".name"] = "Name is required"
		}

		if _, err := transaction.ParsedDate(); err != nil {
			errors[basePath+".date"] = "Date must be in YYYY-MM-DD format"
		}

		if !transaction.Amount.IsPositive() {
			errors[basePath+".amount"] = "Amount must be greater than zero"
		}

		// Income is identified by its category's income flag, so an income row
		// without a category would be counted as spending.
		if transaction.IsIncome && (transaction.CategoryId == nil || *transaction.CategoryId == 0) {
			errors[basePath+".categoryId"] = "An income category is required for a deposit"
		}
	}

	return structs.ValidatorError{Errors: errors}
}
