package analytic

import (
	"context"
	"fmt"
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

// LogDailyMangaView to save manga view event to firebase
func LogDailyMangaView(mangaTitle string, mangaPage int, userIP string) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	analyticsRef := ref.Child("daily_manga_views_data")
	analyticsDailyMangaViewRef := analyticsRef.Child("daily_manga_views_data_val")

	dateNow := fmt.Sprintf(now.Format("2006-01-02"))
	selectedDailyMangaViewRef := analyticsDailyMangaViewRef.Child(dateNow)
	var dailyMangaView models.AnalyticDailyMangaView
	if err := selectedDailyMangaViewRef.Get(ctx, &dailyMangaView); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	dailyMangaView.Title = "MAIN"
	if dailyMangaView.Count == 0 {
		dailyMangaView.ReportDate = dateNow
		dailyMangaView.Count++
		err := selectedDailyMangaViewRef.Set(ctx, &dailyMangaView)
		if err != nil {
			log.Fatalln("Error setting value:", err)
		}
	} else {
		dailyMangaView.Count++
		tempDailyMangaView := map[string]interface{}{
			"report_date": dailyMangaView.ReportDate,
			"title":       dailyMangaView.Title,
			"count":       dailyMangaView.Count,
		}
		selectedDailyMangaViewRef.Update(ctx, tempDailyMangaView)
	}

	selectedDailyMangaViewRefDetailed := selectedDailyMangaViewRef.Child(mangaTitle)
	var dailyMangaViewDetailed models.AnalyticDailyMangaView
	if err := selectedDailyMangaViewRefDetailed.Get(ctx, &dailyMangaViewDetailed); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	dailyMangaViewDetailed.Title = mangaTitle
	if dailyMangaViewDetailed.Count == 0 {
		dailyMangaViewDetailed.ReportDate = dateNow
		dailyMangaViewDetailed.Count++
		err := selectedDailyMangaViewRefDetailed.Set(ctx, &dailyMangaViewDetailed)
		if err != nil {
			log.Fatalln("Error setting value:", err)
		}
	} else {
		dailyMangaViewDetailed.Count++
		tempDailyMangaViewDetailed := map[string]interface{}{
			"report_date": dailyMangaViewDetailed.ReportDate,
			"title":       dailyMangaViewDetailed.Title,
			"count":       dailyMangaViewDetailed.Count,
		}
		selectedDailyMangaViewRefDetailed.Update(ctx, tempDailyMangaViewDetailed)
	}
}
