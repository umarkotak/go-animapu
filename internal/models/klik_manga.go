package models

type KlikMangaSearchParams struct {
	Title  string
	Status string
	Genre  string
}

type KlikMangaHistory struct {
	Title              string `json:"title"`
	CompactTitle       string `json:"compact_title"`
	ImageURL           string `json:"image_url"`
	LastReadChapterInt string `json:"last_read_chapter_int"`
	LastReadChapterID  string `json:"last_read_chapter_id"`
	LastChapterInt     string `json:"last_chapter_int"`
	LastChapterID      string `json:"last_chapter_id"`
	UpdatedAt          string `json:"updated_at"`
	UpdatedAtUnix      string `json:"updated_at_unix"`
}
