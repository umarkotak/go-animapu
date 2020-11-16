package analytic

import (
	"context"
	"log"
	"time"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
)

var ctx = context.Background()

// LogUserAnalytic to save user event to firebase
func LogUserAnalytic(mangaTitle string, mangaPage int, userIP string) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	analyticsRef := ref.Child("analytic_data")
	analyticsUsersRef := analyticsRef.Child("AnalyticDataUsers")

	userAnalyticRef := analyticsUsersRef.Child(userIP)
	var userAnalytic models.AnalyticDataUser
	if err := userAnalyticRef.Get(ctx, &userAnalytic); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	userAnalytic.LastUpdate = now.Format(time.RFC3339)
	userAnalytic.HitCount++
	userAnalytic.UserIP = userIP
	userAnalytic.Location = "Indonesia"

	if userAnalytic.UserIP == "" {
		err := userAnalyticRef.Set(ctx, &userAnalytic)
		if err != nil {
			log.Fatalln("Error setting value:", err)
		}
	} else {
		x := map[string]interface{}{
			"LastUpdate": now.Format(time.RFC3339),
			"HitCount":   userAnalytic.HitCount,
			"Location":   "Indonesia",
		}
		userAnalyticRef.Update(ctx, x)
	}

	mangaAnalyticRef := userAnalyticRef.Child(mangaTitle)
	var mangaAnalytic models.AnalyticDataManga
	if err := mangaAnalyticRef.Get(ctx, &mangaAnalytic); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	mangaAnalytic.HitCount++
	mangaAnalytic.LastUpdate = now.Format(time.RFC3339)
	mangaAnalytic.Page = mangaPage
	mangaAnalytic.Title = mangaTitle

	if mangaAnalytic.Title == "" {
		err := mangaAnalyticRef.Set(ctx, &mangaAnalytic)
		if err != nil {
			log.Fatalln("Error setting value:", err)
		}
	} else {
		x := map[string]interface{}{
			"LastUpdate": now.Format(time.RFC3339),
			"HitCount":   mangaAnalytic.HitCount,
			"Page":       mangaPage,
		}
		mangaAnalyticRef.Update(ctx, x)
	}
}
