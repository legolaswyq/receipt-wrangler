package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
)

type CategoryRepository struct {
	BaseRepository
}

func NewCategoryRepository(tx *gorm.DB) CategoryRepository {
	repository := CategoryRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository CategoryRepository) GetAllCategories(querySelect string) ([]models.Category, error) {
	db := repository.GetDB()
	var categories []models.Category

	err := db.Table("categories").Select(querySelect).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByIds returns the full category records for the given ids (order not guaranteed). Used to
// resolve names for id-only selections (e.g. quick-scan category picks) so they can flow through
// receipt validation, which requires a category name.
func (repository CategoryRepository) GetByIds(ids []uint) ([]models.Category, error) {
	db := repository.GetDB()
	var categories []models.Category

	if len(ids) == 0 {
		return categories, nil
	}

	err := db.Model(&models.Category{}).Where("id IN ?", ids).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// CountByIds returns how many of the given category ids exist. Used to validate
// that a role's category grants reference real categories. Duplicate ids in the
// input are de-duplicated by the IN clause, so callers should pass a unique set.
func (repository CategoryRepository) CountByIds(ids []uint) (int64, error) {
	db := repository.GetDB()

	var count int64
	err := db.Model(&models.Category{}).Where("id IN ?", ids).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository CategoryRepository) CreateCategory(category models.Category) (models.Category, error) {
	db := repository.GetDB()

	// Assign a default palette color when the caller didn't pick one, choosing the
	// first color not already used by another category.
	if len(category.Color) == 0 {
		var usedColors []string
		if err := db.Model(&models.Category{}).
			Where("color <> ''").
			Distinct().
			Pluck("color", &usedColors).Error; err != nil {
			return models.Category{}, err
		}
		category.Color = nextCategoryColor(usedColors)
	}

	err := db.Model(&category).Create(&category).Error
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (repository CategoryRepository) GetAllPagedCategories(pagedRequestCommand commands.PagedRequestCommand) ([]models.CategoryView, error) {
	db := repository.GetDB()
	var categories []models.CategoryView
	quotedAlias := "\"NumberOfReceipts\""

	if pagedRequestCommand.OrderBy == "numberOfReceipts" {
		pagedRequestCommand.OrderBy = quotedAlias
	}

	query := repository.Sort(db, pagedRequestCommand.OrderBy, pagedRequestCommand.SortDirection)
	query = query.Scopes(repository.Paginate(pagedRequestCommand.Page, pagedRequestCommand.PageSize))
	selectString := fmt.Sprintf("categories.id, categories.name, categories.description, categories.is_income, categories.color, COUNT(DISTINCT receipt_categories.receipt_id) as %s", quotedAlias)
	query = query.Table("categories").
		Select(selectString).
		Joins("LEFT JOIN receipt_categories ON categories.id = receipt_categories.category_id").
		Group("categories.id, categories.name")

	err := query.Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (repository CategoryRepository) UpdateCategory(categoryToUpdate models.Category, querySelect string) (models.Category, error) {
	db := repository.GetDB()

	// Map-form Updates writes every listed key regardless of its zero value, unlike GORM's
	// struct-form Updates/Save, which silently skips zero-value fields (e.g. a bool false). That
	// matters here because toggling IsIncome off must persist.
	err := db.Model(models.Category{}).Where("id = ?", categoryToUpdate.ID).Updates(map[string]interface{}{
		"name":        categoryToUpdate.Name,
		"description": categoryToUpdate.Description,
		"is_income":   categoryToUpdate.IsIncome,
		"color":       categoryToUpdate.Color,
	}).Error
	if err != nil {
		return models.Category{}, err
	}

	return categoryToUpdate, nil
}

// GetCategoryById fetches a single category by its id.
func (repository CategoryRepository) GetCategoryById(id string) (models.Category, error) {
	db := repository.GetDB()
	var category models.Category

	err := db.Model(models.Category{}).Where("id = ?", id).First(&category).Error
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (repository CategoryRepository) DeleteCategory(categoryId uint) error {
	db := repository.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&models.ReceiptCategory{}, "category_id = ?", categoryId).Error
		if err != nil {
			return err
		}

		// Remove the category from every group's budget set. The membership delete
		// is raw (no FK cascade), so this must be explicit, matching the other
		// cascade helpers in this codebase.
		budgetRepository := NewCategoryBudgetRepository(tx)
		err = budgetRepository.DeleteBudgetsByCategoryId(tx, categoryId)
		if err != nil {
			return err
		}

		err = tx.Where("id = ?", categoryId).Delete(&models.Category{}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
