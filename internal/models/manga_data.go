package models

// MangaData Detail of certain manga
type MangaData struct {
	CompactTitle     string `json:"compact_title"`
	MangaLastChapter int    `json:"manga_last_chapter"`
	AveragePage      int    `json:"average_page"`
	Status           string `json:"status"`
	ImageURL         string `json:"image_url"`
	NewAdded         int    `json:"new_added"`
	Weight           int    `json:"weight"`
	Finder           string `json:"finder"`
}
