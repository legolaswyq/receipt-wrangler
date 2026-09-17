package repositories

import (
	"receipt-wrangler/api/internal/models"

	"gorm.io/gorm"
)

// Name of the one-time migration that gives categories sharing a color a distinct
// one. The first backfill wrapped the palette once it was exhausted (10 colors),
// so installs with more than 10 categories had duplicates; the palette was later
// widened, and this reassigns the duplicates into the freed colors.
const recolorDuplicateCategoriesMigration = "recolor-duplicate-categories"

// recolorDuplicateCategories walks categories in a stable order, keeps the first
// occurrence of each color, and reassigns any later duplicate (or empty) color to
// the next unused palette color. Distinct colors — including any a user hand-picked
// — are left untouched.
func recolorDuplicateCategories(tx *gorm.DB) error {
	var categories []models.Category
	if err := tx.Order("id asc").Find(&categories).Error; err != nil {
		return err
	}

	usedColors := make([]string, 0, len(categories))
	seen := make(map[string]bool, len(categories))

	for _, category := range categories {
		if category.Color != "" && !seen[category.Color] {
			seen[category.Color] = true
			usedColors = append(usedColors, category.Color)
			continue
		}

		color := nextCategoryColor(usedColors)
		if err := tx.Model(&models.Category{}).
			Where("id = ?", category.ID).
			Update("color", color).Error; err != nil {
			return err
		}
		seen[color] = true
		usedColors = append(usedColors, color)
	}

	return nil
}
