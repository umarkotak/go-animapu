package statistic

import (
	"context"

	"github.com/umarkotak/go-animapu/internal/models"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
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
