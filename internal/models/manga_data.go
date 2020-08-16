package models

// MangaData Detail of certain manga
type MangaData struct {
	MangaLastChapter int    `json:"manga_last_chapter"`
	AveragePage      int    `json:"average_page"`
	Status           string `json:"status"`
	ImageURL         string `json:"image_url"`
	NewAdded         int    `json:"new_added"`
}
