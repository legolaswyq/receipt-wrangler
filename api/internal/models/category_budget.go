package models

import "github.com/shopspring/decimal"

// CategoryBudget is a recurring monthly spend target for a single category
// within a single group. One target per (GroupId, CategoryId).
type CategoryBudget struct {
	BaseModel
	GroupId    uint            `gorm:"not null;uniqueIndex:idx_group_category_budget" json:"groupId"`
	CategoryId uint            `gorm:"not null;uniqueIndex:idx_group_category_budget" json:"categoryId"`
	Amount     decimal.Decimal `gorm:"not null" json:"amount"`
}
