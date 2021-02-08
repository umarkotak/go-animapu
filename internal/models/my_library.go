package models

// MyLibrary contain user handpicked manga
type MyLibrary struct {
	MangaTitle       string `json:"manga_title"`
	MangaLastChapter int    `json:"manga_last_chapter"`
	AveragePage      int    `json:"average_page"`
	Status           string `json:"status"`
	ImageURL         string `json:"image_url"`
	NewAdded         int    `json:"new_added"`
	Weight           int    `json:"weight"`
	Finder           string `json:"finder"`
}
