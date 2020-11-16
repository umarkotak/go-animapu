package models

// AnalyticDataUser user info data
type AnalyticDataUser struct {
	UserIP             string                        `json:"UserIP"`
	HitCount           int                           `json:"HitCount"`
	LastUpdate         string                        `json:"LastUpdate"`
	Location           string                        `json:"Location"`
	AnalyticDataMangas map[string]*AnalyticDataManga `json:"AnalyticDataMangas"`
}
