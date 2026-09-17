package models

import (
	"database/sql/driver"
	"errors"
)

type WidgetType string

const (
	GROUP_SUMMARY     WidgetType = "GROUP_SUMMARY"
	FILTERED_RECEIPTS WidgetType = "FILTERED_RECEIPTS"
	GROUP_ACTIVITY    WidgetType = "GROUP_ACTIVITY"
	PIE_CHART         WidgetType = "PIE_CHART"
	REPORT            WidgetType = "REPORT"
	BUDGET            WidgetType = "BUDGET"
	SPENDING_TABLE    WidgetType = "SPENDING_TABLE"
)

func (widgetType *WidgetType) Scan(value string) error {
	*widgetType = WidgetType(value)
	return nil
}

func (widgetType WidgetType) Value() (driver.Value, error) {
	if widgetType != GROUP_SUMMARY &&
		widgetType != FILTERED_RECEIPTS &&
		widgetType != GROUP_ACTIVITY &&
		widgetType != PIE_CHART &&
		widgetType != REPORT &&
		widgetType != BUDGET &&
		widgetType != SPENDING_TABLE {
		return nil, errors.New("invalid widget type")
	}
	return string(widgetType), nil
}
