package user

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
	rUser "github.com/umarkotak/go-animapu/internal/repository/user"
	sUser "github.com/umarkotak/go-animapu/internal/service/user"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetKlikMangaHistory(c *gin.Context) {
	auth := c.Request.Header["Authorization"][0]

	userData, err := sUser.DetailService(auth)
	if err != nil {
		logrus.Errorf("sUser.DetailService: %v\n", err)
		http_req.RenderResponse(c, 422, err)
	}

	klikMangaHistoriesMap, err := rUser.GetKlikMangaHistories(userData)
	if err != nil {
		logrus.Errorf("rUser.SetKlikMangaHistory: %v\n", err)
		http_req.RenderResponse(c, 422, err)
	}

	klikMangaHistoriesList := []models.KlikMangaHistory{}
	for _, v := range klikMangaHistoriesMap {
		klikMangaHistoriesList = append(klikMangaHistoriesList, v)
	}

	// < ascending
	// > descending
	sort.Slice(klikMangaHistoriesList, func(i, j int) bool {
		return klikMangaHistoriesList[i].UpdatedAtUnix > klikMangaHistoriesList[j].UpdatedAtUnix
	})

	response := map[string]interface{}{
		"klik_manga_histories_map":  klikMangaHistoriesMap,
		"klik_manga_histories_list": klikMangaHistoriesList,
	}

	http_req.RenderResponse(c, 200, response)
}

func PostKlikMangaHistory(c *gin.Context) {
	auth := c.Request.Header["Authorization"][0]

	klikMangaHistory := models.KlikMangaHistory{}

	c.BindJSON(&klikMangaHistory)

	userData, err := sUser.DetailService(auth)
	if err != nil {
		logrus.Errorf("sUser.DetailService: %v\n", err)
		http_req.RenderResponse(c, 422, err)
	}

	err = rUser.SetKlikMangaHistory(userData, klikMangaHistory)
	if err != nil {
		logrus.Errorf("rUser.SetKlikMangaHistory: %v\n", err)
		http_req.RenderResponse(c, 422, err)
	}

	http_req.RenderResponse(c, 200, klikMangaHistory)
}
