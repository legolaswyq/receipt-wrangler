package services

import (
	"errors"
	"fmt"
	"os/exec"
	"receipt-wrangler/api/internal/structs"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

// ErrUnsupportedBankStatement is returned when the uploaded document is not a
// statement layout the parser recognises. It is a user-facing (400) condition,
// not a server fault.
var ErrUnsupportedBankStatement = errors.New("this file is not a supported bank statement. Only ANZ New Zealand transaction PDFs are supported")

const (
	anzBankName            = "ANZ"
	anzStatementDateLayout = "02 Jan 2006"
	anzTransferType        = "Transfer"
)

var (
	anzHeaderRegex = regexp.MustCompile(`^\s*Date\s+Type\s+Details\s+Deposits\s+Withdrawals\s+Balance\s*$`)
	// Date, type, details, one amount (deposit OR withdrawal), running balance.
	// Columns are separated by at least two spaces in pdftotext's layout mode;
	// single spaces occur inside a column ("Visa Purchase", "Pak N Save Mt Albert").
	anzRowRegex     = regexp.MustCompile(`^(\d{2} [A-Z][a-z]{2} \d{4})\s{2,}(\S(?:.*?\S)?)\s{2,}(\S(?:.*?\S)?)\s{2,}(-?\$[\d,]+\.\d{2})\s+(-?\$[\d,]+\.\d{2})\s*$`)
	anzPeriodRegex  = regexp.MustCompile(`Transactions between (\d{2} [A-Z][a-z]{2} \d{4}) and (\d{2} [A-Z][a-z]{2} \d{4})`)
	anzAccountRegex = regexp.MustCompile(`^\s*(\d{2}-\d{4}-\d{7}-\d{2,3})\s*$`)
	anzFooterRegex  = regexp.MustCompile(`Page \d+ of \d+|^Copyright `)
)

// ExtractPdfText returns the text layer of a PDF with its column layout
// preserved, using poppler's pdftotext. Statement parsing depends on that layout
// to tell the deposit column from the withdrawal column.
func ExtractPdfText(pdfPath string) (string, error) {
	output, err := exec.Command("pdftotext", "-layout", pdfPath, "-").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// pdftotext ran and rejected the file: not a readable PDF.
			return "", ErrUnsupportedBankStatement
		}
		return "", fmt.Errorf("pdftotext is not available: %w", err)
	}

	return string(output), nil
}

// ParseAnzStatementText parses the layout-preserved text of an ANZ New Zealand
// transaction PDF. Transactions are returned in statement order (newest first).
func ParseAnzStatementText(text string) (structs.BankStatementPreview, error) {
	preview := structs.BankStatementPreview{
		Bank:         anzBankName,
		Transactions: []structs.BankStatementTransaction{},
	}

	// Each page starts with a form feed glued to its first line.
	lines := strings.Split(strings.ReplaceAll(text, "\f", ""), "\n")

	depositsEnd, withdrawalsEnd := -1, -1
	lastRowIndex := -1

	for i, line := range lines {
		if anzHeaderRegex.MatchString(line) {
			depositsEnd = columnEnd(line, "Deposits")
			withdrawalsEnd = columnEnd(line, "Withdrawals")
			lastRowIndex = -1
			continue
		}

		if depositsEnd < 0 {
			parseAnzPreamble(&preview, lines, i)
			continue
		}

		match := anzRowRegex.FindStringSubmatchIndex(line)
		if match == nil {
			// The line under a row carries its reference (card number, transfer id...).
			if lastRowIndex >= 0 && i == lastRowIndex+1 && isAnzReferenceLine(line) {
				last := &preview.Transactions[len(preview.Transactions)-1]
				last.Reference = strings.TrimSpace(line)
			}
			continue
		}

		transaction, err := buildAnzTransaction(line, match, depositsEnd, withdrawalsEnd)
		if err != nil {
			return preview, err
		}

		preview.Transactions = append(preview.Transactions, transaction)
		lastRowIndex = i
	}

	if depositsEnd < 0 || len(preview.Transactions) == 0 {
		return preview, ErrUnsupportedBankStatement
	}

	err := verifyAnzRunningBalance(preview.Transactions)
	if err != nil {
		return preview, err
	}

	return preview, nil
}

func parseAnzPreamble(preview *structs.BankStatementPreview, lines []string, index int) {
	line := lines[index]

	if period := anzPeriodRegex.FindStringSubmatch(line); period != nil {
		preview.PeriodStart = formatAnzDate(period[1])
		preview.PeriodEnd = formatAnzDate(period[2])
		return
	}

	if account := anzAccountRegex.FindStringSubmatch(line); account != nil {
		preview.AccountNumber = account[1]
		return
	}

	// The very first line is "<account name>   <balance>   <available funds>".
	if index == 0 {
		preview.AccountName = strings.TrimSpace(strings.SplitN(strings.TrimSpace(line), "  ", 2)[0])
	}
}

func buildAnzTransaction(line string, match []int, depositsEnd int, withdrawalsEnd int) (structs.BankStatementTransaction, error) {
	group := func(n int) string { return line[match[2*n]:match[2*n+1]] }

	date, err := time.Parse(anzStatementDateLayout, group(1))
	if err != nil {
		return structs.BankStatementTransaction{}, ErrUnsupportedBankStatement
	}

	amount, err := parseStatementMoney(group(4))
	if err != nil {
		return structs.BankStatementTransaction{}, ErrUnsupportedBankStatement
	}

	balance, err := parseStatementMoney(group(5))
	if err != nil {
		return structs.BankStatementTransaction{}, ErrUnsupportedBankStatement
	}

	// Amounts are right-aligned under their column header, so whichever header
	// ends closer to where the amount ends is the column it sits in.
	amountEnd := utf8.RuneCountInString(line[:match[9]])
	isIncome := absInt(amountEnd-depositsEnd) < absInt(amountEnd-withdrawalsEnd)

	transactionType := group(2)
	details := group(3)

	return structs.BankStatementTransaction{
		Date:     date.Format(structs.BankStatementDateFormat),
		Type:     transactionType,
		Details:  details,
		Name:     buildStatementReceiptName(transactionType, details),
		Amount:   amount,
		Balance:  balance,
		IsIncome: isIncome,
		Selected: transactionType != anzTransferType,
	}, nil
}

// verifyAnzRunningBalance cross-checks the column-based deposit/withdrawal
// decision against the statement's own running balance. Rows are newest first,
// so each row's balance must equal the next (older) row's balance plus or minus
// its amount. A mismatch means the layout was misread, and importing would
// silently turn spending into income (or vice versa) — fail instead.
func verifyAnzRunningBalance(transactions []structs.BankStatementTransaction) error {
	for i := 0; i < len(transactions)-1; i++ {
		current := transactions[i]
		signedAmount := current.Amount
		if !current.IsIncome {
			signedAmount = signedAmount.Neg()
		}

		if !transactions[i+1].Balance.Add(signedAmount).Equal(current.Balance) {
			return fmt.Errorf("%w (running balance does not add up at %s %s)",
				ErrUnsupportedBankStatement, current.Date, current.Details)
		}
	}

	return nil
}

func isAnzReferenceLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return len(trimmed) > 0 && strings.HasPrefix(line, " ") && !anzFooterRegex.MatchString(trimmed)
}

func buildStatementReceiptName(transactionType string, details string) string {
	if transactionType == anzTransferType {
		return transactionType + " " + details
	}

	return details
}

func columnEnd(line string, header string) int {
	index := strings.Index(line, header)
	if index < 0 {
		return -1
	}

	return utf8.RuneCountInString(line[:index+len(header)])
}

func parseStatementMoney(value string) (decimal.Decimal, error) {
	cleaned := strings.NewReplacer("$", "", ",", "").Replace(value)
	return decimal.NewFromString(cleaned)
}

func formatAnzDate(value string) string {
	date, err := time.Parse(anzStatementDateLayout, value)
	if err != nil {
		return ""
	}

	return date.Format(structs.BankStatementDateFormat)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
