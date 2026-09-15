package repositories_test

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCategoryBudgetUpsertAndList(t *testing.T) {
	defer repositories.TruncateTestDb()
	repo := repositories.NewCategoryBudgetRepository(nil)

	_, err := repo.UpsertBudget(1, 5, decimal.NewFromInt(400))
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	// Upsert again updates in place, does not duplicate.
	_, err = repo.UpsertBudget(1, 5, decimal.NewFromInt(450))
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	budgets, err := repo.GetBudgetsByGroupId(1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(budgets) != 1 {
		t.Fatalf("expected 1 budget, got %d", len(budgets))
	}
	if !budgets[0].Amount.Equal(decimal.NewFromInt(450)) {
		t.Errorf("expected amount 450, got %s", budgets[0].Amount.String())
	}

	if err := repo.DeleteBudget(1, 5); err != nil {
		t.Fatalf("delete: %v", err)
	}
	budgets, _ = repo.GetBudgetsByGroupId(1)
	if len(budgets) != 0 {
		t.Errorf("expected 0 budgets after delete, got %d", len(budgets))
	}
}

func TestDeleteCategoryRemovesBudgets(t *testing.T) {
	defer repositories.TruncateTestDb()
	catRepo := repositories.NewCategoryRepository(nil)
	budgetRepo := repositories.NewCategoryBudgetRepository(nil)

	cat, err := catRepo.CreateCategory(models.Category{Name: "Dining Out"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := budgetRepo.UpsertBudget(1, cat.ID, decimal.NewFromInt(400)); err != nil {
		t.Fatalf("upsert budget: %v", err)
	}

	if err := catRepo.DeleteCategory(cat.ID); err != nil {
		t.Fatalf("delete category: %v", err)
	}

	budgets, _ := budgetRepo.GetBudgetsByGroupId(1)
	if len(budgets) != 0 {
		t.Errorf("expected budgets removed with category, got %d", len(budgets))
	}
}
