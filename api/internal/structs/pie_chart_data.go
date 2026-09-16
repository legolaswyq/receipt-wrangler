package structs

type PieChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color,omitempty"`
}

type PieChartData struct {
	Data []PieChartDataPoint `json:"data"`
}
