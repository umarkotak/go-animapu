package manga

import (
	"context"
	"log"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
)

var ctx = context.Background()

// SetUserToFirebase set user to firebase
func SetUserToFirebase(userData models.UserData) models.UserData {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	userRef := ref.Child("user_db")
	err := userRef.Set(ctx, map[string]*models.UserData{
		"hello": &userData,
	})

	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return userData
}

// GetUserByUsernameFromFirebase get user from firebase by username
func GetUserByUsernameFromFirebase(username string) models.UserData {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("user_db")
	userRef := ref.Child(username)
	var userData models.UserData
	if err := userRef.Get(ctx, &userData); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	return userData
}
