package models

type MangaChapterDetail struct {
	Title       string   `json:"title"`
	Chapter     int      `json:"chapter"`
	NextChapter int      `json:"next_chapter"`
	PrevChapter int      `json:"prev_chapter"`
	Images      []string `json:"images"`
}
