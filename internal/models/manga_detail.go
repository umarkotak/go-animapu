package models

type MangaDetail struct {
	Title          string       `json:"title"`
	CompactTitle   string       `json:"compact_title"`
	Chapters       []string     `json:"chapters"`
	ChaptersInt    []int64      `json:"chapters_int"`
	ChapterLinks   []string     `json:"chapter_links"`
	Genres         string       `json:"genres"`
	ImageURL       string       `json:"image_url"`
	Description    string       `json:"description"`
	MangadexID     string       `json:"mangadex_id"`
	MangaUpdatesID string       `json:"manga_updates_id"`
	LastChapter    string       `json:"last_chapter"`
	LastChapterInt int64        `json:"last_chapter_int"`
	DetailLink     string       `json:"detail_link"`
	ChapterObjs    []ChapterObj `json:"chapter_objs"`
}
