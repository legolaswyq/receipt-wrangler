package structs

import "github.com/shopspring/decimal"

// BankStatementDateFormat is the wire format of every date on a bank statement
// preview/import payload. Statement dates are calendar dates with no time or
// zone, so they travel as plain YYYY-MM-DD strings.
const BankStatementDateFormat = "2006-01-02"

// BankStatementTransaction is one parsed row of a bank statement, offered to the
// user for review before any receipt is created.
type BankStatementTransaction struct {
	// AccountName and AccountNumber identify the statement the row came from, so
	// rows from several accounts can be reviewed in one list.
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`

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
	// Transfer is true when the row moves money between the user's own accounts,
	// which is neither spending nor income: the bank labelled it a transfer, or
	// its other side was found in another uploaded statement.
	Transfer bool `json:"transfer"`
	// Selected is the suggested default for the review step: false for transfers
	// between the user's own accounts and for rows that look already imported.
	Selected            bool  `json:"selected"`
	Duplicate           bool  `json:"duplicate"`
	SuggestedCategoryId *uint `json:"suggestedCategoryId"`
}

// BankStatementAccount describes one uploaded statement.
type BankStatementAccount struct {
	Bank          string `json:"bank"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	PeriodStart   string `json:"periodStart"`
	PeriodEnd     string `json:"periodEnd"`
}

// BankStatementPreview is the review payload for one or more statements: the
// accounts they cover, and every transaction across them, newest first.
type BankStatementPreview struct {
	Accounts     []BankStatementAccount     `json:"accounts"`
	Transactions []BankStatementTransaction `json:"transactions"`
}

type BankStatementImportResult struct {
	CreatedCount int `json:"createdCount"`
}
