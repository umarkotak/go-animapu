package firebase_helper

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

var (
	ctx         = context.Background()
	firebaseApp *firebase.App
)

// GetFirebaseApp return firebase app instance
func GetFirebaseApp() *firebase.App {
	if firebaseApp != nil {
		return firebaseApp
	}

	conf := &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL"),
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_GOOGLE_APPLICATION_CREDENTIALS"))

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	firebaseApp = app
	return firebaseApp
}

// GetFirebaseDB db instance
func GetFirebaseDB() *db.Client {
	client, err := firebaseApp.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client
}
