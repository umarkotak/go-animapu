package models

// AnalyticData overall analytic data
type AnalyticData struct {
	TotalHitCount     int                          `json:"TotalHitCount"`
	AnalyticDataUsers map[string]*AnalyticDataUser `json:"AnalyticDataUsers"`
}
