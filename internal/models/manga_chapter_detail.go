package models

type MangaChapterDetail struct {
	Title          string   `json:"title"`
	MangadexID     string   `json:"mangadex_id"`
	MangaUpdatesID string   `json:"manga_updates_id"`
	Chapter        int      `json:"chapter"`
	NextChapter    int      `json:"next_chapter"`
	PrevChapter    int      `json:"prev_chapter"`
	Images         []string `json:"images"`
}
