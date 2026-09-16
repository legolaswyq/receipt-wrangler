package models

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/utils"
)

type Category struct {
	BaseModel
	Name        string `gorm:"not null; uniqueIndex" json:"name"`
	Description string `json:"description"`
	IsIncome    bool   `gorm:"not null;default:false" json:"isIncome"`
	Color       string `json:"color"`
}

func (category *Category) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &category)
	if err != nil {
		return err
	}

	return nil
}
