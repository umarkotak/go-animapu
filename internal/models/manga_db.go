package models

// MangaDB manga database representative
type MangaDB struct {
	MangaDatas map[string]*MangaData `json:"manga_db"`
}
