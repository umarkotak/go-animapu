package manga

import (
	"github.com/gin-gonic/gin"
	sStatistic "github.com/umarkotak/go-animapu/internal/service/statistic"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetMangaStatistics(c *gin.Context) {
	result := sStatistic.GenerateMangaStatistic()

	http_req.RenderResponse(c, 200, result)
}

func GetDailyMangaStatistics(c *gin.Context) {
	result := sStatistic.GenerateDailyMangaStatistic()

	http_req.RenderResponse(c, 200, result)
}
