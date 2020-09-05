package models

// MangaDB manga database representative
type MangaDB struct {
	MangaDatas map[string]*MangaData `json:"manga_db"`
}

// MangaDBFromInterface return MangaDB from map[string]interface{}
func MangaDBFromInterface(mangaData map[string]interface{}) MangaDB {
	mangaDB := MangaDB{
		MangaDatas: make(map[string]*MangaData),
	}

	for key, manga := range mangaData {
		mappedManga := manga.(map[string]interface{})
		dataMangaLastChapter := int(mappedManga["manga_last_chapter"].(float64))
		dataAveragePage := int(mappedManga["average_page"].(float64))
		dataStatus := mappedManga["status"].(string)
		dataImageURL := mappedManga["image_url"].(string)
		dataNewAdded := int(mappedManga["new_added"].(float64))
		dataWeight := int(mappedManga["weight"].(float64))

		mangaData := &MangaData{
			MangaLastChapter: dataMangaLastChapter,
			AveragePage:      dataAveragePage,
			Status:           dataStatus,
			ImageURL:         dataImageURL,
			NewAdded:         dataNewAdded,
			Weight:           dataWeight,
		}

		mangaDB.MangaDatas[key] = mangaData
	}

	return mangaDB
}
