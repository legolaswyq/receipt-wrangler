package repositories

import (
	"receipt-wrangler/api/internal/models"

	"gorm.io/gorm"
)

// Name of the one-time migration that assigns a default palette color to every
// existing category that has none. New categories get a color at creation time
// (CreateCategory); this covers rows that predate the color field.
const backfillCategoryColorsMigration = "backfill-category-colors"

// backfillCategoryColors walks categories with an empty color in a stable order
// and assigns each the next unused palette color, so an upgraded install gets
// distinct, stable pie-chart colors without any manual step.
func backfillCategoryColors(tx *gorm.DB) error {
	var categories []models.Category
	if err := tx.Where("color IS NULL OR color = ''").
		Order("id asc").
		Find(&categories).Error; err != nil {
		return err
	}

	var usedColors []string
	if err := tx.Model(&models.Category{}).
		Where("color <> ''").
		Distinct().
		Pluck("color", &usedColors).Error; err != nil {
		return err
	}

	for _, category := range categories {
		color := nextCategoryColor(usedColors)
		if err := tx.Model(&models.Category{}).
			Where("id = ?", category.ID).
			Update("color", color).Error; err != nil {
			return err
		}
		usedColors = append(usedColors, color)
	}

	return nil
}
