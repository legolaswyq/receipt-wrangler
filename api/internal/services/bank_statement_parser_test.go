package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// statementRow is one synthetic transaction for buildAnzStatementText. Exactly
// one of deposit/withdrawal is set.
type statementRow struct {
	date       string
	kind       string
	details    string
	deposit    string
	withdrawal string
	balance    string
	reference  string
}

// buildAnzStatementText lays rows out the way pdftotext -layout prints an ANZ
// statement: amounts right-aligned under their column header, the reference on
// the following line, and a new page (form feed + repeated header, at different
// column offsets) every pageSize rows.
func buildAnzStatementText(rows []statementRow, pageSize int) string {
	var b strings.Builder
	b.WriteString("Everyday Account                                              $100.00             $100.00\n")
	b.WriteString("MS TEST PERSON                                            Account Balance     Available Funds\n")
	b.WriteString("Go\n")
	b.WriteString("01-0001-0000001-00\n")
	b.WriteString("                                              Transactions between 01 Jan 2026 and 31 Jan 2026\n")

	for i, row := range rows {
		page := i / pageSize
		// Later pages are narrower, as on a real statement.
		detailsWidth := 50 - 12*page
		if i%pageSize == 0 {
			if page > 0 {
				fmt.Fprintf(&b, "\n\n01 Feb 2026 09:00 NZT%sPage %d of 9\n", strings.Repeat(" ", 60), page)
				b.WriteString("Copyright 2026 © ANZ Bank New Zealand Limited. All rights reserved.\n\f")
			}
			fmt.Fprintf(&b, "%-14s%-19s%-*s%12s%16s%12s\n", "Date", "Type", detailsWidth, "Details", "Deposits", "Withdrawals", "Balance")
		}

		fmt.Fprintf(&b, "%-14s%-19s%-*s%12s%16s%12s\n", row.date, row.kind, detailsWidth, row.details, row.deposit, row.withdrawal, row.balance)
		fmt.Fprintf(&b, "%s%s\n", strings.Repeat(" ", 33), row.reference)
	}

	return b.String()
}

func sampleStatementRows() []statementRow {
	// Newest first; each balance is the next row's balance plus/minus the amount.
	return []statementRow{
		{date: "20 Jan 2026", kind: "Visa Purchase", details: "Corner Groc", withdrawal: "$1,042.19", balance: "$1,317.81", reference: "Card number: 1111 **** **** 2222"},
		{date: "19 Jan 2026", kind: "Payment", details: "Acme Trading", withdrawal: "$40.00", balance: "$2,360.00"},
		{date: "18 Jan 2026", kind: "Transfer", details: "From: 06-0002-0000002-01", deposit: "$1,000.00", balance: "$2,400.00", reference: "Credit Transfer 123456"},
		{date: "15 Jan 2026", kind: "Direct Credit", details: "Employer Limited", deposit: "$1,500.00", balance: "$1,400.00", reference: "Salary"},
		{date: "02 Jan 2026", kind: "EFTPOS", details: "Pak N Save Mt Albert", withdrawal: "$100.00", balance: "-$100.00", reference: "Card number: 1111 **** **** 2222"},
	}
}

func TestParseAnzStatementText_ParsesRowsAcrossPages(t *testing.T) {
	preview, err := ParseAnzStatementText(buildAnzStatementText(sampleStatementRows(), 2))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if preview.Bank != "ANZ" || preview.AccountName != "Everyday Account" || preview.AccountNumber != "01-0001-0000001-00" {
		t.Errorf("unexpected header: %+v", preview)
	}
	if preview.PeriodStart != "2026-01-01" || preview.PeriodEnd != "2026-01-31" {
		t.Errorf("unexpected period: %s to %s", preview.PeriodStart, preview.PeriodEnd)
	}
	if len(preview.Transactions) != 5 {
		t.Fatalf("expected 5 transactions, got %d", len(preview.Transactions))
	}

	expected := []struct {
		date, kind, name, amount, balance, reference string
		isIncome, selected                           bool
	}{
		{"2026-01-20", "Visa Purchase", "Corner Groc", "1042.19", "1317.81", "Card number: 1111 **** **** 2222", false, true},
		{"2026-01-19", "Payment", "Acme Trading", "40", "2360", "", false, true},
		{"2026-01-18", "Transfer", "Transfer From: 06-0002-0000002-01", "1000", "2400", "Credit Transfer 123456", true, false},
		{"2026-01-15", "Direct Credit", "Employer Limited", "1500", "1400", "Salary", true, true},
		{"2026-01-02", "EFTPOS", "Pak N Save Mt Albert", "100", "-100", "Card number: 1111 **** **** 2222", false, true},
	}

	for i, want := range expected {
		got := preview.Transactions[i]
		if got.Date != want.date || got.Type != want.kind || got.Name != want.name || got.Reference != want.reference {
			t.Errorf("row %d: unexpected text fields: %+v", i, got)
		}
		if got.Amount.String() != want.amount || got.Balance.String() != want.balance {
			t.Errorf("row %d: expected amount %s balance %s, got %s / %s", i, want.amount, want.balance, got.Amount, got.Balance)
		}
		if got.IsIncome != want.isIncome || got.Selected != want.selected {
			t.Errorf("row %d: expected isIncome=%t selected=%t, got %t / %t", i, want.isIncome, want.selected, got.IsIncome, got.Selected)
		}
		if got.Transfer != (want.kind == "Transfer") {
			t.Errorf("row %d: expected only bank-labelled transfers to be flagged, got transfer=%t", i, got.Transfer)
		}
		if got.AccountName != "Everyday Account" || got.AccountNumber != "01-0001-0000001-00" {
			t.Errorf("row %d: expected the statement's account on the row, got %q %q", i, got.AccountName, got.AccountNumber)
		}
		if got.Duplicate || got.SuggestedCategoryId != nil {
			t.Errorf("row %d: parser must not set duplicate/suggestion", i)
		}
	}
}

func TestParseAnzStatementText_RejectsUnrecognisedText(t *testing.T) {
	tests := map[string]string{
		"empty":            "",
		"not a statement":  "Invoice 42\nTotal due   $10.00\n",
		"header, no rows":  "Date          Type               Details            Deposits     Withdrawals     Balance\n",
		"rows, no header":  "20 Jan 2026   Visa Purchase      Corner Groc                         $42.19   $1,317.81\n",
		"only a page foot": "01 Feb 2026 09:00 NZT                Page 1 of 1\n",
	}

	for name, text := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ParseAnzStatementText(text)
			if !errors.Is(err, ErrUnsupportedBankStatement) {
				t.Errorf("expected ErrUnsupportedBankStatement, got %v", err)
			}
		})
	}
}

// A statement whose running balance contradicts the column an amount was read
// from must be refused: importing it would swap income and spending.
func TestParseAnzStatementText_RejectsBalanceMismatch(t *testing.T) {
	rows := sampleStatementRows()
	// Move the salary into the withdrawals column; its balance still says deposit.
	rows[3].withdrawal, rows[3].deposit = rows[3].deposit, ""

	_, err := ParseAnzStatementText(buildAnzStatementText(rows, 10))
	if !errors.Is(err, ErrUnsupportedBankStatement) {
		t.Fatalf("expected ErrUnsupportedBankStatement, got %v", err)
	}
	if !strings.Contains(err.Error(), "Employer Limited") {
		t.Errorf("expected the error to name the offending row, got %q", err.Error())
	}
}

func TestExtractPdfText_RejectsNonPdf(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext is not installed")
	}

	path := filepath.Join(t.TempDir(), "not-a.pdf")
	if err := os.WriteFile(path, []byte("this is not a pdf"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := ExtractPdfText(path)
	if !errors.Is(err, ErrUnsupportedBankStatement) {
		t.Errorf("expected ErrUnsupportedBankStatement, got %v", err)
	}
}
