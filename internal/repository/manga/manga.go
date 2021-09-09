package manga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	pkgAppCache "github.com/umarkotak/go-animapu/internal/lib/app_cache"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/lib/firebase_helper"
	"github.com/umarkotak/go-animapu/internal/models"
)

var mangaDbFilePath = "data/mangas.json"
var ctx = context.Background()

// var mangaDbFilePath = "data/mangas.json"

// GetMangaFromJSON read manga from json db in data/
func GetMangaFromJSON() models.MangaDB {
	mangaDbJSONFile, err := os.Open(mangaDbFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer mangaDbJSONFile.Close()

	mangaDbByteValue, _ := ioutil.ReadAll(mangaDbJSONFile)

	var mangaDB models.MangaDB
	json.Unmarshal([]byte(mangaDbByteValue), &mangaDB)

	return mangaDB
}

// UpdateMangaToJSON update manga to json db in data/
func UpdateMangaToJSON(mangaDb models.MangaDB) models.MangaDB {
	mangaDbJSON, _ := json.Marshal(mangaDb)
	ioutil.WriteFile(mangaDbFilePath, mangaDbJSON, 0644)

	return mangaDb
}

// GetMangaFromFireBaseV1 fetch manga from firebase using map string interface
func GetMangaFromFireBaseV1() models.MangaDB {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("manga_db")
	var mangaDBDataFirebase map[string]interface{}
	if err := ref.Get(ctx, &mangaDBDataFirebase); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	mangaDB := models.MangaDBFromInterface(mangaDBDataFirebase)

	return mangaDB
}

// GetMangaFromFireBaseV2 using direct struct
func GetMangaFromFireBaseV2() models.MangaDB {
	appCache := pkgAppCache.GetAppCache()

	res, found := appCache.Get("manga_from_firebase_v2")
	if found {
		fmt.Println("FETCH FROM APP CACHE")
		return res.(models.MangaDB)
	}

	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	var mangaDB models.MangaDB
	if err := ref.Get(ctx, &mangaDB); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	appCache.Set("manga_from_firebase_v2", mangaDB, 5*time.Minute)
	return mangaDB
}

// GetMangaFromFireBaseV2WithoutCache using direct struct
func GetMangaFromFireBaseV2WithoutCache() models.MangaDB {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	var mangaDB models.MangaDB
	if err := ref.Get(ctx, &mangaDB); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	return mangaDB
}

// UpdateMangaToFireBase update manga data to firebase
func UpdateMangaToFireBase(mangaDB models.MangaDB) models.MangaDB {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	mangaRef := ref.Child("manga_db")
	err := mangaRef.Set(ctx, mangaDB.MangaDatas)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return mangaDB
}

func AddMangaToFireBaseGeneralLibrary(mangaData models.MangaData) (models.MangaData, error) {
	firebaseDB := firebaseHelper.GetFirebaseDB()
	rootRef := firebaseDB.NewRef("")
	mangaDBRef := rootRef.Child("manga_db")

	dataRef := mangaDBRef.Child(mangaData.Title)

	var tempData models.MangaData
	dataRef.Get(ctx, &tempData)
	if tempData.MangaLastChapter != 0 {
		return tempData, errors.New("already_exists")
	}

	err := dataRef.Set(ctx, mangaData)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return mangaData, nil
}
