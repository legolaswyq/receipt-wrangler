package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func createBudgetTestReceipt(
	name string,
	amount float64,
	paidByUserId uint,
	groupId uint,
	date time.Time,
	categories []models.Category,
) models.Receipt {
	db := repositories.GetDB()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromFloat(amount),
		Date:         date,
		PaidByUserID: paidByUserId,
		GroupId:      groupId,
		Status:       models.OPEN,
		Categories:   categories,
	}
	db.Create(&receipt)
	return receipt
}

func TestGetBudgetData_SplitsIncomeSpendUntracked(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()

	db := repositories.GetDB()

	salaryCategory := models.Category{Name: "Salary", IsIncome: true}
	diningCategory := models.Category{Name: "Dining Out"}
	gasCategory := models.Category{Name: "Gas"}
	db.Create(&salaryCategory)
	db.Create(&diningCategory)
	db.Create(&gasCategory)

	budgetRepository := repositories.NewCategoryBudgetRepository(nil)
	_, err := budgetRepository.UpsertBudget(1, diningCategory.ID, decimal.NewFromInt(400))
	if err != nil {
		t.Fatalf("UpsertBudget: %v", err)
	}

	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)
	thisMonth := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	lastMonth := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)

	// Receipts this month.
	createBudgetTestReceipt("Paycheck", 5000, 1, 1, thisMonth, []models.Category{salaryCategory})
	createBudgetTestReceipt("Dinner", 300, 1, 1, thisMonth, []models.Category{diningCategory})
	createBudgetTestReceipt("Gas Station", 90, 1, 1, thisMonth, []models.Category{gasCategory})

	// Receipt last month, must be excluded.
	createBudgetTestReceipt("Old Dinner", 999, 1, 1, lastMonth, []models.Category{diningCategory})

	svc := NewBudgetService(nil)
	data, err := svc.GetBudgetData(1, "1", now)
	if err != nil {
		t.Fatalf("GetBudgetData: %v", err)
	}

	if data.Income != 5000 {
		t.Errorf("income = %v, want 5000", data.Income)
	}
	if data.Spent != 390 { // 300 + 90, last-month 999 excluded
		t.Errorf("spent = %v, want 390", data.Spent)
	}
	if data.Untracked != 90 { // Gas has no budget
		t.Errorf("untracked = %v, want 90", data.Untracked)
	}
	if data.Net != 4610 { // 5000 - 390
		t.Errorf("net = %v, want 4610", data.Net)
	}
	if len(data.Categories) != 1 || data.Categories[0].Name != "Dining Out" {
		t.Fatalf("expected 1 budgeted category Dining Out, got %+v", data.Categories)
	}
	if data.Categories[0].Spent != 300 || data.Categories[0].Target != 400 || data.Categories[0].Over {
		t.Errorf("dining bar wrong: %+v", data.Categories[0])
	}
}
