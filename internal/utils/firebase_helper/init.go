package firebase_helper

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

var ctx = context.Background()

// GetFirebaseApp return firebase app instance
func GetFirebaseApp() *firebase.App {
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

	return app
}

// GetFirebaseDB db instance
func GetFirebaseDB() *db.Client {
	client, err := GetFirebaseApp().Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client
}
