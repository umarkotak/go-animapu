package models

// ReadHistory user manga read history
type ReadHistory struct {
	MangaTitle     string `json:"manga_title"`
	LastChapter    int    `json:"last_chapter"`
	LastReadTime   string `json:"last_read_time"`
	LastReadTimeI  int64  `json:"last_read_time_i"`
	MangadexID     string `json:"mangadex_id"`
	MangaUpdatesID string `json:"manga_updates_id"`
	MangaSource    string `json:"manga_source"`
}
