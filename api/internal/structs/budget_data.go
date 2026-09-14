package structs

type BudgetCategory struct {
	CategoryId uint    `json:"categoryId"`
	Name       string  `json:"name"`
	Target     float64 `json:"target"`
	Spent      float64 `json:"spent"`
	Over       bool    `json:"over"`
}

type BudgetData struct {
	Month      string           `json:"month"`
	Income     float64          `json:"income"`
	Spent      float64          `json:"spent"`
	Net        float64          `json:"net"`
	Untracked  float64          `json:"untracked"`
	Categories []BudgetCategory `json:"categories"`
}
