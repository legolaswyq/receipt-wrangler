package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// seedStatementImport creates a user who is a member of a group, plus one
// expense and one income category.
func seedStatementImport(t *testing.T) (models.User, models.Group, models.Category, models.Category) {
	t.Helper()
	db := repositories.GetDB()

	user := models.User{Username: "statement-user", Password: "p", DisplayName: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	group := models.Group{Name: "statement-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: user.ID}).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	groceries, salary := seedStatementCategories(t)
	return user, group, groceries, salary
}

func seedStatementCategories(t *testing.T) (models.Category, models.Category) {
	t.Helper()
	db := repositories.GetDB()

	groceries := models.Category{Name: "Groceries"}
	if err := db.Create(&groceries).Error; err != nil {
		t.Fatalf("seed expense category: %v", err)
	}
	salary := models.Category{Name: "Salary", IsIncome: true}
	if err := db.Create(&salary).Error; err != nil {
		t.Fatalf("seed income category: %v", err)
	}

	return groceries, salary
}

func uintPointer(value uint) *uint {
	return &value
}

func TestBankStatementPreview_FlagsDuplicatesAndDefaultsSoleIncomeCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, _, salary := seedStatementImport(t)
	db := repositories.GetDB()

	// One of the two identical-looking candidates already exists as a receipt.
	existing := models.Receipt{
		Name:         "acme trading",
		Amount:       decimal.RequireFromString("40.00"),
		Date:         time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC),
		GroupId:      group.ID,
		PaidByUserID: user.ID,
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed existing receipt: %v", err)
	}

	service := NewBankStatementImportService(nil)
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{buildAnzStatementText(sampleStatementRows(), 10)})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}

	for _, transaction := range preview.Transactions {
		isExisting := transaction.Name == "Acme Trading"
		if transaction.Duplicate != isExisting {
			t.Errorf("%s: expected duplicate=%t, got %t", transaction.Name, isExisting, transaction.Duplicate)
		}
		if isExisting && transaction.Selected {
			t.Errorf("a duplicate must not be selected by default")
		}

		// No AI is configured: deposits fall back to the only income category and
		// withdrawals get no suggestion.
		switch {
		case transaction.IsIncome && (transaction.SuggestedCategoryId == nil || *transaction.SuggestedCategoryId != salary.ID):
			t.Errorf("%s: expected the sole income category to be suggested", transaction.Name)
		case !transaction.IsIncome && transaction.SuggestedCategoryId != nil:
			t.Errorf("%s: expected no suggestion for a withdrawal without AI", transaction.Name)
		}
	}
}

func TestBankStatementPreview_DuplicateMatchingConsumesOneReceiptPerRow(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, _, _ := seedStatementImport(t)
	db := repositories.GetDB()

	rows := []statementRow{
		{date: "05 Jan 2026", kind: "Visa Purchase", details: "7 Mart", withdrawal: "$2.50", balance: "$95.00"},
		{date: "05 Jan 2026", kind: "Visa Purchase", details: "7 Mart", withdrawal: "$2.50", balance: "$97.50"},
	}
	existing := models.Receipt{
		Name:         "7 Mart",
		Amount:       decimal.RequireFromString("2.5"),
		Date:         time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		GroupId:      group.ID,
		PaidByUserID: user.ID,
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed existing receipt: %v", err)
	}

	service := NewBankStatementImportService(nil)
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{buildAnzStatementText(rows, 10)})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}

	duplicates := 0
	for _, transaction := range preview.Transactions {
		if transaction.Duplicate {
			duplicates++
		}
	}
	if duplicates != 1 {
		t.Errorf("expected exactly one of two identical rows to be a duplicate, got %d", duplicates)
	}
}

func TestBankStatementPreview_AppliesAiSuggestionsOfTheRightKind(t *testing.T) {
	defer repositories.TruncateTestDb()
	body := ""
	server := newMutableOllamaServer(t, &body)
	user, group, _ := seedReceiptImagePipeline(t, server.URL)
	groceries, salary := seedStatementCategories(t)
	db := repositories.GetDB()
	bonus := models.Category{Name: "Bonus", IsIncome: true}
	if err := db.Create(&bonus).Error; err != nil {
		t.Fatalf("seed second income category: %v", err)
	}

	// Merchant indexes follow first appearance: 0 Corner Groc, 1 Acme Trading,
	// 2 the transfer, 3 Employer Limited, 4 Pak N Save. The model answers one
	// correctly, puts an expense category on a deposit, an income category on a
	// withdrawal, and invents an id.
	body = ollamaReceiptResponse(t, `{"suggestions":[
		{"index":0,"categoryId":`+uintString(groceries.ID)+`},
		{"index":1,"categoryId":`+uintString(salary.ID)+`},
		{"index":2,"categoryId":`+uintString(groceries.ID)+`},
		{"index":3,"categoryId":`+uintString(bonus.ID)+`},
		{"index":4,"categoryId":99999}
	]}`)

	service := NewBankStatementImportService(nil)
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{buildAnzStatementText(sampleStatementRows(), 10)})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}

	expected := map[string]*uint{
		"Corner Groc":                       &groceries.ID,
		"Acme Trading":                      nil,
		"Transfer From: 06-0002-0000002-01": nil,
		"Employer Limited":                  &bonus.ID,
		"Pak N Save Mt Albert":              nil,
	}
	for _, transaction := range preview.Transactions {
		want := expected[transaction.Name]
		got := transaction.SuggestedCategoryId
		if (want == nil) != (got == nil) || (want != nil && *want != *got) {
			t.Errorf("%s: expected suggestion %v, got %v", transaction.Name, formatUintPointer(want), formatUintPointer(got))
		}
	}
}

// buildSavingsStatementText is a second account's statement for the same
// period, so previews can be exercised across accounts.
func buildSavingsStatementText(rows []statementRow) string {
	text := buildAnzStatementText(rows, 10)
	text = strings.Replace(text, "Everyday Account", "Savings Account ", 1)
	return strings.Replace(text, "01-0001-0000001-00", "01-0001-0000001-50", 1)
}

func TestBankStatementPreview_PairsTransfersAcrossUploadedAccounts(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, _, _ := seedStatementImport(t)

	everyday := buildAnzStatementText([]statementRow{
		{date: "20 Jan 2026", kind: "Automatic Payment", details: "T Person", withdrawal: "$20.00", balance: "$1,440.00", reference: "To Savings Acc"},
		{date: "20 Jan 2026", kind: "Payment", details: "Acme Trading", withdrawal: "$40.00", balance: "$1,460.00"},
		{date: "15 Jan 2026", kind: "Direct Credit", details: "Employer Limited", deposit: "$1,500.00", balance: "$1,500.00", reference: "Salary"},
	}, 10)
	savings := buildSavingsStatementText([]statementRow{
		{date: "20 Jan 2026", kind: "Credit", details: "T Person", deposit: "$20.00", balance: "$60.00", reference: "From Everyday Acc"},
		// Same day and amount as the Acme payment, but also a withdrawal: not its other side.
		{date: "20 Jan 2026", kind: "Payment", details: "Other Shop", withdrawal: "$40.00", balance: "$40.00"},
		// Same day and amount as the Acme payment and a deposit, but from someone
		// else: an unrelated $40, not the other side of that payment.
		{date: "20 Jan 2026", kind: "Direct Credit", details: "A Friend", deposit: "$40.00", balance: "$80.00"},
		{date: "19 Jan 2026", kind: "Direct Credit", details: "Interest Paid", deposit: "$40.00", balance: "$40.00"},
	})

	service := NewBankStatementImportService(nil)
	// The everyday statement is picked twice; its rows must appear once.
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{everyday, savings, everyday})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}

	if len(preview.Accounts) != 2 || preview.Accounts[0].AccountNumber != "01-0001-0000001-00" || preview.Accounts[1].AccountNumber != "01-0001-0000001-50" {
		t.Fatalf("expected the two distinct accounts, got %+v", preview.Accounts)
	}
	if len(preview.Transactions) != 7 {
		t.Fatalf("expected 7 transactions, got %d", len(preview.Transactions))
	}

	for i, transaction := range preview.Transactions {
		if i > 0 && transaction.Date > preview.Transactions[i-1].Date {
			t.Errorf("rows must be newest first across accounts, got %s after %s", transaction.Date, preview.Transactions[i-1].Date)
		}

		isTransfer := transaction.Name == "T Person"
		if transaction.Transfer != isTransfer {
			t.Errorf("%s %s: expected transfer=%t, got %t", transaction.AccountName, transaction.Name, isTransfer, transaction.Transfer)
		}
		if transaction.Selected == isTransfer {
			t.Errorf("%s %s: expected selected=%t", transaction.AccountName, transaction.Name, !isTransfer)
		}
	}
}

func TestBankStatementPreview_DoesNotPairWithinOneAccount(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, _, _ := seedStatementImport(t)

	// A purchase refunded the same day: same account, so both rows are real.
	text := buildAnzStatementText([]statementRow{
		{date: "20 Jan 2026", kind: "Credit", details: "Corner Groc", deposit: "$30.00", balance: "$100.00"},
		{date: "20 Jan 2026", kind: "Visa Purchase", details: "Corner Groc", withdrawal: "$30.00", balance: "$70.00"},
	}, 10)

	service := NewBankStatementImportService(nil)
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{text})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}

	for _, transaction := range preview.Transactions {
		if transaction.Transfer || !transaction.Selected {
			t.Errorf("%s: rows of a single account must not be paired as a transfer", transaction.Name)
		}
	}
}

func TestBankStatementPreview_SurvivesAiFailure(t *testing.T) {
	defer repositories.TruncateTestDb()
	body := "not json at all"
	server := newMutableOllamaServer(t, &body)
	user, group, _ := seedReceiptImagePipeline(t, server.URL)
	seedStatementCategories(t)

	service := NewBankStatementImportService(nil)
	preview, err := service.previewFromTexts(user.ID, group.ID, []string{buildAnzStatementText(sampleStatementRows(), 10)})
	if err != nil {
		t.Fatalf("a failed suggestion call must not fail the preview: %v", err)
	}
	if len(preview.Transactions) != 5 {
		t.Errorf("expected 5 transactions, got %d", len(preview.Transactions))
	}
}

func TestBankStatementImport_CreatesResolvedReceiptsWithCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, groceries, salary := seedStatementImport(t)
	db := repositories.GetDB()

	command := commands.BankStatementImportCommand{
		GroupId:      group.ID,
		PaidByUserId: user.ID,
		Transactions: []commands.BankStatementImportTransaction{
			{Date: "2026-01-20", Name: " Corner Groc ", Amount: decimal.RequireFromString("42.19"), CategoryId: uintPointer(groceries.ID)},
			{Date: "2026-01-15", Name: "Employer Limited", Amount: decimal.RequireFromString("1500"), IsIncome: true, CategoryId: uintPointer(salary.ID)},
			{Date: "2026-01-02", Name: "Uncategorised Shop", Amount: decimal.RequireFromString("100")},
		},
	}

	service := NewBankStatementImportService(nil)
	result, err := service.Import(user.ID, command)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if result.CreatedCount != 3 {
		t.Errorf("expected 3 created, got %d", result.CreatedCount)
	}

	var receipts []models.Receipt
	if err := db.Preload("Categories").Order("date desc").Find(&receipts).Error; err != nil {
		t.Fatalf("load receipts: %v", err)
	}
	if len(receipts) != 3 {
		t.Fatalf("expected 3 receipts, got %d", len(receipts))
	}

	first := receipts[0]
	if first.Name != "Corner Groc" || !first.Amount.Equal(decimal.RequireFromString("42.19")) {
		t.Errorf("unexpected first receipt: %s %s", first.Name, first.Amount)
	}
	if !first.Date.Equal(time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected the date at UTC midnight, got %s", first.Date)
	}
	if first.Status != models.RESOLVED || first.ResolvedDate == nil {
		t.Errorf("expected a resolved receipt with a resolved date, got %s / %v", first.Status, first.ResolvedDate)
	}
	if first.GroupId != group.ID || first.PaidByUserID != user.ID || first.CreatedBy == nil || *first.CreatedBy != user.ID {
		t.Errorf("unexpected ownership on first receipt: %+v", first)
	}
	if len(first.Categories) != 1 || first.Categories[0].ID != groceries.ID {
		t.Errorf("expected the groceries category, got %+v", first.Categories)
	}
	if len(receipts[1].Categories) != 1 || !receipts[1].Categories[0].IsIncome {
		t.Errorf("expected the income category on the deposit, got %+v", receipts[1].Categories)
	}
	if len(receipts[2].Categories) != 0 {
		t.Errorf("expected no category on the uncategorised row, got %+v", receipts[2].Categories)
	}

	// Attaching a category must not rewrite the category row itself.
	var reloaded models.Category
	db.First(&reloaded, groceries.ID)
	if reloaded.Name != "Groceries" {
		t.Errorf("category name was changed to %q", reloaded.Name)
	}
}

func TestBankStatementImport_RejectsInvalidSelectionsAndCreatesNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, groceries, salary := seedStatementImport(t)
	db := repositories.GetDB()
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	if err := db.Create(&allGroup).Error; err != nil {
		t.Fatalf("seed all group: %v", err)
	}

	valid := commands.BankStatementImportTransaction{Date: "2026-01-20", Name: "Corner Groc", Amount: decimal.RequireFromString("42.19"), CategoryId: uintPointer(groceries.ID)}

	tests := map[string]commands.BankStatementImportCommand{
		"expense category on a deposit": {GroupId: group.ID, PaidByUserId: user.ID, Transactions: []commands.BankStatementImportTransaction{
			valid,
			{Date: "2026-01-15", Name: "Employer Limited", Amount: decimal.RequireFromString("1500"), IsIncome: true, CategoryId: uintPointer(groceries.ID)},
		}},
		"income category on a withdrawal": {GroupId: group.ID, PaidByUserId: user.ID, Transactions: []commands.BankStatementImportTransaction{
			valid,
			{Date: "2026-01-15", Name: "Shop", Amount: decimal.RequireFromString("15"), CategoryId: uintPointer(salary.ID)},
		}},
		"unknown category": {GroupId: group.ID, PaidByUserId: user.ID, Transactions: []commands.BankStatementImportTransaction{
			valid,
			{Date: "2026-01-15", Name: "Shop", Amount: decimal.RequireFromString("15"), CategoryId: uintPointer(99999)},
		}},
		"all group": {GroupId: allGroup.ID, PaidByUserId: user.ID, Transactions: []commands.BankStatementImportTransaction{valid}},
	}

	service := NewBankStatementImportService(nil)
	for name, command := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := service.Import(user.ID, command)
			if !errors.Is(err, ErrBankStatementInvalidSelection) {
				t.Errorf("expected ErrBankStatementInvalidSelection, got %v", err)
			}

			var count int64
			db.Model(&models.Receipt{}).Count(&count)
			if count != 0 {
				t.Errorf("expected no receipts to be created, found %d", count)
			}
		})
	}
}

func uintString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}

func formatUintPointer(value *uint) string {
	if value == nil {
		return "nil"
	}

	return uintString(*value)
}
