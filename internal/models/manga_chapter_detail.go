package models

type MangaChapterDetail struct {
	Title          string       `json:"title"`
	MangadexID     string       `json:"mangadex_id"`
	MangaUpdatesID string       `json:"manga_updates_id"`
	Chapter        int          `json:"chapter"`
	ChapterObjs    []ChapterObj `json:"chapter_objs"`
	NextChapter    int          `json:"next_chapter"`
	PrevChapter    int          `json:"prev_chapter"`
	Images         []string     `json:"images"`
}

type ChapterObj struct {
	Title  string  `json:"title"`
	Link   string  `json:"link"`
	Number float64 `json:"number"`
}
