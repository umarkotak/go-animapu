package models

// AnalyticDailyMangaView daily manga view
type AnalyticDailyMangaView struct {
	ReportDate string `json:"report_date"`
	Title      string `json:"title"`
	Count      int    `json:"count"`
}
