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
	firebaseHelper "github.com/umarkotak/go-animapu/internal/utils/firebase_helper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

var (
	port string
)

func Init() {
	godotenv.Load(".env")
	pkgAppCache.InitAppCache()

	firebaseHelper.GetFirebaseApp()

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
	router.POST("/mangas/general/add", cManga.PostAddToGeneralMangaLibrary)

	router.GET("/mangas/statistics", cManga.GetMangaStatistics)
	router.GET("/mangas/daily_manga_statistics", cManga.GetDailyMangaStatistics)

	router.GET("/mangas_detail", cManga.GetMangaDetail)
	router.GET("/mangas/search_v1", cManga.GetMangaSearch)
	router.GET("/mangas/todays_v1", cManga.GetMangaTodays)

	router.GET("/mangas/maid_my/home", cManga.GetMaidMyHome)
	router.GET("/mangas/maid_my/search", cManga.GetMaidMySearch)
	router.GET("/mangas/maid_my/manga_detail", cManga.GetMaidMyMangaDetail)
	router.GET("/mangas/maid_my/manga_chapter_detail", cManga.GetMaidMyMangaChapterDetail)

	router.GET("/mangas/klik_manga/home/:page", cManga.GetKlikMangaHomeNextPage)
	router.GET("/mangas/klik_manga/home", cManga.GetKlikMangaHome)
	router.GET("/mangas/klik_manga/search/:page", cManga.GetKlikMangaSearchNextPage)
	router.GET("/mangas/klik_manga/search", cManga.GetKlikMangaSearch)
	router.GET("/mangas/klik_manga/manga/:manga_id", cManga.GetKlikMangaDetail)
	router.GET("/mangas/klik_manga/manga/:manga_id/:manga_chapter", cManga.GetKlikMangaChapterDetail)

	// users
	router.POST("/users/register", cUser.RegisterUserFirebase)
	router.POST("/users/login", cUser.LoginUser)
	router.POST("/users/read_histories", cUser.LogReadHistories)
	router.GET("/users/detail", cUser.GetDetailFirebase)
	router.POST("/users/add_manga_library", cUser.AddToMyMangaLibrary)
	router.GET("/users/manga_library", cUser.GetMyLibrary)
	router.POST("/users/remove_manga_library", cUser.RemoveMyLibrary)
	router.POST("/users/analytic_v1", cAnalytic.PostUserAnalyticV1)
	router.POST("/users/klik_manga/history", cUser.PostKlikMangaHistory)
	router.GET("/users/klik_manga/history", cUser.GetKlikMangaHistory)

	// clip
	router.GET("/clips", cClips.GetClips)
	router.POST("/clips", cClips.CreateClip)

	// animepahe
	router.GET("/animes_map", cAnimes.GetAnimesMap)
	router.GET("/search_anime", cAnimes.GetSearchAnime)

	// mangadex proxy
	router.GET("/mangadex/*mangadex_path", cMangadex.GetProxy)

	// mangaupdates scrapper
	router.GET("/mangaupdates/releases/:mangaupdates_id", cMangaupdates.ReleasesSearch)
	router.GET("/mangaupdates/releases", cMangaupdates.GetReleases)
	router.GET("/mangaupdates/series", cMangaupdates.GetSeries)
	router.GET("/mangaupdates/search", cMangaupdates.Search)
	router.GET("/mangaupdates/detail/:manga_title", cMangaupdates.GetDetailByTitle)

	// klik manga scrapper
	// HOME: https://klikmanga.com/
	// DETAIL WITH CHAPTER: https://klikmanga.com/manga/solo-leveling/
	// CHAPTER: https://klikmanga.com/manga/solo-leveling/chapter-167/
	// SEARCH: https://klikmanga.com/?s=solo+leveling&post_type=wp-manga

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
	http_req.RenderResponse(c, 500, gin.H{
		"success": false,
		"data":    nil,
		"error":   fmt.Sprintf("Internal server error: %v", err),
	})
	return
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			http_req.RenderResponse(c, 200, gin.H{
				"success": true,
				"data":    nil,
				"error":   "",
			})
			return
		}
		c.Next()
	}
}
