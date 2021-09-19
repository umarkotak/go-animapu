package user

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/utils/firebase_helper"
)

var ctx = context.Background()

// SetUserToFirebase set user to firebase
func SetUserToFirebase(userData models.UserData) models.UserData {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	userRef := ref.Child("user_db")

	if userData.Username == "" {
		return userData
	}

	userDataRef := userRef.Child(userData.Username)
	err := userDataRef.Set(ctx, &userData)

	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return userData
}

// GetUserByUsernameFromFirebase get user from firebase by username
func GetUserByUsernameFromFirebase(username string) models.UserData {
	// appCache := pkgAppCache.GetAppCache()

	// res, found := appCache.Get("user_firebase")
	// if found {
	// 	fmt.Println("FETCH FROM APP CACHE")
	// 	return res.(models.UserData)
	// }

	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("user_db")
	userRef := ref.Child(username)
	var userData models.UserData
	if err := userRef.Get(ctx, &userData); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	// appCache.Set("user_firebase", userData, 5*time.Minute)

	return userData
}

func SetMangaToMyLibrary(userData models.UserData, myLibrary models.MyLibrary) {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("user_db")

	if userData.Username == "" {
		return
	}

	userRef := ref.Child(userData.Username)
	myLibraryRef := userRef.Child("MyLibraries")
	selectedMangaRef := myLibraryRef.Child(myLibrary.MangaTitle)
	err := selectedMangaRef.Set(ctx, &myLibrary)
	if err != nil {
		log.Fatalln("Error set to database:", err)
	}
}

func RemoveMangaFromMyLibrary(userData models.UserData, myLibrary models.MyLibrary) {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("user_db")

	if userData.Username == "" {
		return
	}

	userRef := ref.Child(userData.Username)
	myLibraryRef := userRef.Child("MyLibraries")

	if myLibrary.MangaTitle == "" {
		return
	}

	selectedMangaRef := myLibraryRef.Child(myLibrary.MangaTitle)
	selectedMangaRef.Delete(ctx)
}

func SetKlikMangaHistory(userData models.UserData, klikMangaHistory models.KlikMangaHistory) error {
	if userData.Username == "" {
		return fmt.Errorf("User not found")
	}

	if klikMangaHistory.Title == "" {
		return fmt.Errorf("Title cannot be empty")
	}

	firebaseDB := firebaseHelper.GetFirebaseDB()

	userRef := firebaseDB.NewRef("user_db")
	selectedUserRef := userRef.Child(userData.Username)
	klikMangaHistoriesRef := selectedUserRef.Child("klik_manga_histories_map")
	selectedKlikMangaHistoriesRef := klikMangaHistoriesRef.Child(klikMangaHistory.Title)
	err := selectedKlikMangaHistoriesRef.Set(context.Background(), klikMangaHistory)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

func GetKlikMangaHistories(userData models.UserData) (map[string]models.KlikMangaHistory, error) {
	klikMangaHistoriesMap := map[string]models.KlikMangaHistory{}

	if userData.Username == "" {
		return klikMangaHistoriesMap, fmt.Errorf("User not found")
	}

	firebaseDB := firebaseHelper.GetFirebaseDB()

	userRef := firebaseDB.NewRef("user_db")
	selectedUserRef := userRef.Child(userData.Username)
	klikMangaHistoriesRef := selectedUserRef.Child("klik_manga_histories_map")
	err := klikMangaHistoriesRef.Get(context.Background(), &klikMangaHistoriesMap)
	if err != nil {
		logrus.Errorln(err)
		return klikMangaHistoriesMap, err
	}

	return klikMangaHistoriesMap, nil
}
