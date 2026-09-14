package repositories

import (
	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryBudgetRepository struct {
	BaseRepository
}

func NewCategoryBudgetRepository(tx *gorm.DB) CategoryBudgetRepository {
	return CategoryBudgetRepository{BaseRepository: BaseRepository{DB: GetDB(), TX: tx}}
}

func (repo CategoryBudgetRepository) GetBudgetsByGroupId(groupId uint) ([]models.CategoryBudget, error) {
	db := repo.GetDB()
	var budgets []models.CategoryBudget
	err := db.Where("group_id = ?", groupId).Find(&budgets).Error
	return budgets, err
}

func (repo CategoryBudgetRepository) UpsertBudget(groupId, categoryId uint, amount decimal.Decimal) (models.CategoryBudget, error) {
	db := repo.GetDB()
	budget := models.CategoryBudget{GroupId: groupId, CategoryId: categoryId, Amount: amount}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}, {Name: "category_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount", "updated_at"}),
	}).Create(&budget).Error
	return budget, err
}

func (repo CategoryBudgetRepository) DeleteBudget(groupId, categoryId uint) error {
	db := repo.GetDB()
	return db.Where("group_id = ? AND category_id = ?", groupId, categoryId).
		Delete(&models.CategoryBudget{}).Error
}

func (repo CategoryBudgetRepository) DeleteBudgetsByCategoryId(tx *gorm.DB, categoryId uint) error {
	db := repo.DB
	if tx != nil {
		db = tx
	}
	return db.Where("category_id = ?", categoryId).Delete(&models.CategoryBudget{}).Error
}

func (repo CategoryBudgetRepository) DeleteBudgetsByGroupId(tx *gorm.DB, groupId uint) error {
	db := repo.DB
	if tx != nil {
		db = tx
	}
	return db.Where("group_id = ?", groupId).Delete(&models.CategoryBudget{}).Error
}
