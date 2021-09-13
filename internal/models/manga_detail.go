package models

type MangaDetail struct {
	Title          string   `json:"title"`
	Chapters       []string `json:"chapters"`
	ChapterLinks   []string `json:"chapter_links"`
	Genres         string   `json:"genres"`
	ImageURL       string   `json:"image_url"`
	Description    string   `json:"description"`
	MangadexID     string   `json:"mangadex_id"`
	MangaUpdatesID string   `json:"manga_updates_id"`
}
