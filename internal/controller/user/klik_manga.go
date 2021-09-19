package user

import (
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

	http_req.RenderResponse(c, 200, klikMangaHistoriesMap)
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
