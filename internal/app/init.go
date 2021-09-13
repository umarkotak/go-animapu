package app

import (
	"fmt"
	"os"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	baseController "github.com/umarkotak/go-animapu/internal/controller"
	cAnalytic "github.com/umarkotak/go-animapu/internal/controller/analytic"
	cAnimes "github.com/umarkotak/go-animapu/internal/controller/animes"
	cChats "github.com/umarkotak/go-animapu/internal/controller/chats"
	cClips "github.com/umarkotak/go-animapu/internal/controller/clips"
	cManga "github.com/umarkotak/go-animapu/internal/controller/manga"
	cMangadex "github.com/umarkotak/go-animapu/internal/controller/mangadex"
	cMangaupdates "github.com/umarkotak/go-animapu/internal/controller/mangaupdates"
	cSocketGame "github.com/umarkotak/go-animapu/internal/controller/socket_game"
	cUser "github.com/umarkotak/go-animapu/internal/controller/user"
	"github.com/umarkotak/go-animapu/internal/models"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/utils/app_cache"
)

var (
	port string
)

func Init() {
	godotenv.Load(".env")
	pkgAppCache.InitAppCache()

	port = os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
}

// Start this is entry porint for all http request to go-animapu web
func Start() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(nice.Recovery(recoveryHandler))

	router.GET("/ping", baseController.GetPing)

	// mangas
	router.GET("/mangas", cManga.GetManga)
	router.GET("/mangas/update", cManga.UpdateManga)
	router.GET("/mangas/firebase", cManga.GetMangaFirebase)
	router.GET("/mangas/firebase/update", cManga.UpdateMangaFirebase)
	router.GET("/mangas/search_v1", cManga.GetMangaSearch)
	router.GET("/mangas/todays_v1", cManga.GetMangaTodays)
	router.GET("/mangas/statistics", cManga.GetMangaStatistics)
	router.GET("/mangas/daily_manga_statistics", cManga.GetDailyMangaStatistics)
	router.GET("/mangas_detail", cManga.GetMangaDetail)
	router.GET("/mangas/maid_my/home", cManga.GetMaidMyHome)
	router.GET("/mangas/maid_my/search", cManga.GetMaidMySearch)
	router.GET("/mangas/maid_my/manga_detail", cManga.GetMaidMyMangaDetail)
	router.GET("/mangas/maid_my/manga_chapter_detail", cManga.GetMaidMyMangaChapterDetail)
	router.POST("/mangas/general/add", cManga.PostAddToGeneralMangaLibrary)

	// users
	router.POST("/users/register", cUser.RegisterUserFirebase)
	router.POST("/users/login", cUser.LoginUser)
	router.POST("/users/read_histories", cUser.LogReadHistories)
	router.GET("/users/detail", cUser.GetDetailFirebase)
	router.POST("/users/add_manga_library", cUser.AddToMyMangaLibrary)
	router.GET("/users/manga_library", cUser.GetMyLibrary)
	router.POST("/users/remove_manga_library", cUser.RemoveMyLibrary)

	// users analytic
	router.POST("/users/analytic_v1", cAnalytic.PostUserAnalyticV1)

	// clip
	router.GET("/clips", cClips.GetClips)
	router.POST("/clips", cClips.CreateClip)

	// animepahe
	router.GET("/animes_map", cAnimes.GetAnimesMap)
	router.GET("/search_anime", cAnimes.GetSearchAnime)

	// mangadex proxy
	router.GET("/mangadex/*mangadex_path", cMangadex.GetProxy)

	// mangaupdates scrapper
	// /mangaupdates/releases only -> https://www.mangaupdates.com/releases.html
	// /mangaupdates/search -> https://www.mangaupdates.com/search.html?search=naruto
	// /mangaupdates/series -> https://www.mangaupdates.com/series.html?id=19230
	router.GET("/mangaupdates/series", cMangaupdates.GetSeries)
	// /mangaupdates/releases with query -> https://www.mangaupdates.com/releases.html?stype=series&search=19230&page=1&perpage=100&orderby=chap&asc=desc

	// web sockets

	hub := models.NewHub()
	go hub.Run()

	router.GET("/chats_v1", func(c *gin.Context) {
		cChats.GetChatsConnection(c, hub)
	})

	world := models.NewWorld()
	go world.Run()

	router.GET("/socket_game_v1", func(c *gin.Context) {
		cSocketGame.Serve(c, world)
	})

	router.Use(CORSMiddleware())
	router.Run(":" + port)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(
		500,
		gin.H{
			"success": false,
			"data":    nil,
			"error":   fmt.Sprintf("Internal server error: %v", err),
		},
	)
	c.AbortWithStatus(500)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.JSON(
				200,
				gin.H{
					"success": true,
					"data":    nil,
					"error":   "",
				},
			)
			return
		}

		c.Next()
	}
}
