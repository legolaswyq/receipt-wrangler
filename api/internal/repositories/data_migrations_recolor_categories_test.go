package repositories

import (
	"receipt-wrangler/api/internal/models"
	"testing"
)

func TestRecolorDuplicateCategories(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	// Three categories sharing one color, one with a distinct (hand-picked) color.
	dupA := models.Category{Name: "A", Color: "#4E79A7"}
	dupB := models.Category{Name: "B", Color: "#4E79A7"}
	dupC := models.Category{Name: "C", Color: "#4E79A7"}
	distinct := models.Category{Name: "D", Color: "#123456"}
	db.Create(&dupA)
	db.Create(&dupB)
	db.Create(&dupC)
	db.Create(&distinct)

	if err := recolorDuplicateCategories(db); err != nil {
		t.Fatalf("recolor: %v", err)
	}

	var categories []models.Category
	db.Order("id asc").Find(&categories)

	colors := make(map[string]int)
	for _, c := range categories {
		colors[c.Color]++
	}
	for color, count := range colors {
		if count != 1 {
			t.Errorf("color %s used %d times, want all distinct", color, count)
		}
	}

	// The hand-picked distinct color must be preserved.
	var reloaded models.Category
	db.First(&reloaded, distinct.ID)
	if reloaded.Color != "#123456" {
		t.Errorf("hand-picked color changed to %s", reloaded.Color)
	}
}
