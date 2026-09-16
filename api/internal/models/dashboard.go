package models

type Dashboard struct {
	BaseModel
	Name    string   `gorm:"not null;" json:"name"`
	User    User     `json:"-"`
	UserID  uint     `gorm:"not null;" json:"userId"`
	Group   Group    `json:"-"`
	GroupID uint     `json:"groupId"`
	Widgets []Widget `json:"widgets"`
	// Period is the dashboard-level date range applied to spend-over-time widgets
	// (e.g. the pie chart). It holds a preset key (THIS_MONTH, LAST_MONTH,
	// LAST_3_MONTHS, THIS_YEAR, ALL_TIME, CUSTOM); PeriodStartDate/PeriodEndDate
	// carry the ISO dates only when Period is CUSTOM. Empty means "not set" — the
	// client defaults it to THIS_MONTH.
	Period          string `json:"period"`
	PeriodStartDate string `json:"periodStartDate"`
	PeriodEndDate   string `json:"periodEndDate"`
}
