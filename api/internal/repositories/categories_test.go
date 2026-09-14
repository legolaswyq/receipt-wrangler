package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpdateCategoryPersistsIsIncome(t *testing.T) {
	defer TruncateTestDb()
	repo := NewCategoryRepository(nil)

	created, err := repo.CreateCategory(models.Category{Name: "Salary"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	created.IsIncome = true
	updated, err := repo.UpdateCategory(created, "name, description, is_income")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !updated.IsIncome {
		t.Errorf("expected IsIncome true after update, got false")
	}

	fetched, err := repo.GetCategoryById(utils.UintToString(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !fetched.IsIncome {
		t.Errorf("expected IsIncome persisted, got false")
	}
}
