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

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/utils/app_cache"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/utils/firebase_helper"
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

func UpdateOneMangaToFireBase(mangaData models.MangaData) models.MangaData {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	mangaDBRef := ref.Child("manga_db")
	selectedMangaRef := mangaDBRef.Child(mangaData.Title)
	err := selectedMangaRef.Set(ctx, mangaData)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return mangaData
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

func SetMangaDBToCache(key string, mangaDB models.MangaDB) error {
	logrus.Infof("SetMangaDBToCache, Key: %v\n", key)

	mangaDBByte, err := json.Marshal(mangaDB)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	pkgAppCache.GetAppCache().Set(key, string(mangaDBByte), 5*time.Minute)
	return nil
}

func GetMangaDBFromCache(key string) (models.MangaDB, error) {
	logrus.Infof("GetMangaDBFromCache, Key: %v\n", key)

	var err error
	mangaDB := models.MangaDB{}
	intf, found := pkgAppCache.GetAppCache().Get(key)
	if !found || intf == nil {
		err = fmt.Errorf("Cache Not Found, Key: %v\n", key)
		logrus.Errorln(err)
		return mangaDB, err
	}
	err = json.Unmarshal([]byte(intf.(string)), &mangaDB)
	if err != nil {
		logrus.Errorln(err)
		return mangaDB, err
	}
	return mangaDB, nil
}

func SetObjectToCache(key string, object interface{}) error {
	logrus.Infof("SetMangaDBToCache, Key: %v\n", key)

	mangaDBByte, err := json.Marshal(object)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	pkgAppCache.GetAppCache().Set(key, string(mangaDBByte), 5*time.Minute)
	return nil
}

func GetAnimapuMangasFromCache(key string, destinationObject *[]models.MangaData) error {
	logrus.Infof("GetMangaDBFromCache, Key: %v\n", key)

	var err error
	intf, found := pkgAppCache.GetAppCache().Get(key)
	if !found || intf == nil {
		err = fmt.Errorf("Cache Not Found, Key: %v\n", key)
		logrus.Errorln(err)
		return err
	}
	err = json.Unmarshal([]byte(intf.(string)), destinationObject)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
