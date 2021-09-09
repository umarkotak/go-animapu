package statistic

import (
	"context"
	"log"

	firebaseHelper "github.com/umarkotak/go-animapu/internal/lib/firebase_helper"
	"github.com/umarkotak/go-animapu/internal/models"
)

var ctx = context.Background()

type TempStatisticResult struct {
	Title         string
	TotalHitCount int
}

// GenerateMangaStatistic generate user read manga statistics
func GenerateMangaStatistic() map[string]TempStatisticResult {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	analyticsRef := ref.Child("analytic_data")
	analyticsDataUserRef := analyticsRef.Child("AnalyticDataUsers")
	var mangaAnalytics map[string]map[string]models.AnalyticDataManga
	analyticsDataUserRef.Get(ctx, &mangaAnalytics)

	tempStatisticFinalResult := make(map[string]TempStatisticResult)

	for _, mangaAnalytic := range mangaAnalytics {
		for mangaTitle, mangaAnalyticDetail := range mangaAnalytic {
			if mangaAnalyticDetail.HitCount == 0 {
				continue
			}
			_, ok := tempStatisticFinalResult[mangaTitle]

			if ok {
				tempStatisticFinalResult[mangaTitle] = TempStatisticResult{
					Title:         mangaTitle,
					TotalHitCount: tempStatisticFinalResult[mangaTitle].TotalHitCount + mangaAnalyticDetail.HitCount,
				}
			} else {
				tempStatisticFinalResult[mangaTitle] = TempStatisticResult{
					Title:         mangaTitle,
					TotalHitCount: mangaAnalyticDetail.HitCount,
				}
			}
		}
	}

	return tempStatisticFinalResult
}

// GenerateDailyMangaStatistic generate daily manga statistics
func GenerateDailyMangaStatistic() map[string]models.AnalyticDailyMangaView {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	analyticsRef := ref.Child("daily_manga_views_data")
	analyticsDailyMangaViewRef := analyticsRef.Child("daily_manga_views_data_val")
	var tempResult map[string]models.AnalyticDailyMangaView
	if err := analyticsDailyMangaViewRef.Get(ctx, &tempResult); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	return tempResult
}
