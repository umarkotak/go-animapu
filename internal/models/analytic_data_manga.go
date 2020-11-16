package models

// AnalyticDataManga manga read by the user
type AnalyticDataManga struct {
	Title      string `json:"Title"`
	HitCount   int    `json:"HitCount"`
	Page       int    `json:"Page"`
	LastUpdate string `json:"LastUpdate"`
}
