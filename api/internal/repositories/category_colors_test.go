package repositories_test

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func TestCreateCategoryAssignsDefaultColorFromPalette(t *testing.T) {
	defer repositories.TruncateTestDb()
	repo := repositories.NewCategoryRepository(nil)

	first, err := repo.CreateCategory(models.Category{Name: "Groceries"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if first.Color != repositories.CategoryColorPalette[0] {
		t.Errorf("first category color = %q, want %q", first.Color, repositories.CategoryColorPalette[0])
	}

	second, err := repo.CreateCategory(models.Category{Name: "Dining Out"})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if second.Color != repositories.CategoryColorPalette[1] {
		t.Errorf("second category color = %q, want %q (next unused)", second.Color, repositories.CategoryColorPalette[1])
	}
}

func TestCreateCategoryKeepsExplicitColor(t *testing.T) {
	defer repositories.TruncateTestDb()
	repo := repositories.NewCategoryRepository(nil)

	created, err := repo.CreateCategory(models.Category{Name: "Custom", Color: "#123456"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.Color != "#123456" {
		t.Errorf("explicit color not kept: got %q", created.Color)
	}
}
