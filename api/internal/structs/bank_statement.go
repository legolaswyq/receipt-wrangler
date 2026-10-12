package structs

import "github.com/shopspring/decimal"

// BankStatementDateFormat is the wire format of every date on a bank statement
// preview/import payload. Statement dates are calendar dates with no time or
// zone, so they travel as plain YYYY-MM-DD strings.
const BankStatementDateFormat = "2006-01-02"

// BankStatementTransaction is one parsed row of a bank statement, offered to the
// user for review before any receipt is created.
type BankStatementTransaction struct {
	Date      string          `json:"date"`
	Type      string          `json:"type"`
	Details   string          `json:"details"`
	Reference string          `json:"reference"`
	Name      string          `json:"name"`
	Amount    decimal.Decimal `json:"amount"`
	Balance   decimal.Decimal `json:"balance"`
	// IsIncome is true for a deposit and false for a withdrawal. Amount is always
	// positive; income is identified by this flag, matching Category.IsIncome.
	IsIncome bool `json:"isIncome"`
	// Selected is the suggested default for the review step: false for transfers
	// between the user's own accounts and for rows that look already imported.
	Selected            bool  `json:"selected"`
	Duplicate           bool  `json:"duplicate"`
	SuggestedCategoryId *uint `json:"suggestedCategoryId"`
}

type BankStatementPreview struct {
	Bank          string                     `json:"bank"`
	AccountName   string                     `json:"accountName"`
	AccountNumber string                     `json:"accountNumber"`
	PeriodStart   string                     `json:"periodStart"`
	PeriodEnd     string                     `json:"periodEnd"`
	Transactions  []BankStatementTransaction `json:"transactions"`
}

type BankStatementImportResult struct {
	CreatedCount int `json:"createdCount"`
}
