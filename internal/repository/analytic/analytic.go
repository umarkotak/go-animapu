package analytic

import (
	"context"
	"log"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
)

var ctx = context.Background()

func UpdateAnalyticToFireBase(params models.AnalyticData) models.AnalyticData {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	analyticDataManga := models.AnalyticDataManga{
		Title:      "Mahouka koukou",
		HitCount:   1,
		Page:       1,
		LastUpdate: "2020-01-01",
	}
	analyticDataUser := models.AnalyticDataUser{
		UserIP:     "192-168-1-1",
		Location:   "Indonesia",
		HitCount:   100,
		LastUpdate: "2020-01-01",
		AnalyticDataMangas: map[string]*models.AnalyticDataManga{
			"maohuka koukou": &analyticDataManga,
		},
	}
	analyticData := models.AnalyticData{
		TotalHitCount: 100,
		AnalyticDataUsers: map[string]*models.AnalyticDataUser{
			"192-168-1-1": &analyticDataUser,
		},
	}

	ref := firebaseDB.NewRef("")
	analyticsRef := ref.Child("analytic_data")
	err := analyticsRef.Set(ctx, analyticData)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	return analyticData
}
