package models

// ReadHistory user manga read history
type ReadHistory struct {
	MangaTitle    string `json:"manga_title"`
	LastChapter   int    `json:"last_chapter"`
	LastReadTime  string `json:"last_read_time"`
	LastReadTimeI int64  `json:"last_read_time_i"`
}
