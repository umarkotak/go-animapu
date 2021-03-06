package network

import (
	"github.com/gin-gonic/gin"

	baseController "github.com/umarkotak/go-animapu/internal/controller"
	cAnalytic "github.com/umarkotak/go-animapu/internal/controller/analytic"
	cAnimes "github.com/umarkotak/go-animapu/internal/controller/animes"
	cChats "github.com/umarkotak/go-animapu/internal/controller/chats"
	cClips "github.com/umarkotak/go-animapu/internal/controller/clips"
	cManga "github.com/umarkotak/go-animapu/internal/controller/manga"
	cSocketGame "github.com/umarkotak/go-animapu/internal/controller/socket_game"
	cUser "github.com/umarkotak/go-animapu/internal/controller/user"
	"github.com/umarkotak/go-animapu/internal/models"
)

// RouterStart this is entry porint for all http request to go-animapu web
func RouterStart(port string) {
	router := gin.New()
	router.Use(gin.Logger())

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

	// users
	router.POST("/users/register", cUser.RegisterUserFirebase)
	router.OPTIONS("/users/register", baseController.SkipCors)
	router.POST("/users/login", cUser.LoginUser)
	router.OPTIONS("/users/login", baseController.SkipCors)
	router.POST("/users/read_histories", cUser.LogReadHistories)
	router.OPTIONS("/users/read_histories", baseController.SkipCors)
	router.GET("/users/detail", cUser.GetDetailFirebase)
	router.OPTIONS("/users/detail", baseController.SkipCors)
	router.POST("/users/add_manga_library", cUser.AddToMyMangaLibrary)
	router.OPTIONS("/users/add_manga_library", baseController.SkipCors)
	router.GET("/users/manga_library", cUser.GetMyLibrary)
	router.OPTIONS("/users/manga_library", baseController.SkipCors)
	router.POST("/users/remove_manga_library", cUser.RemoveMyLibrary)
	router.OPTIONS("/users/remove_manga_library", baseController.SkipCors)

	// users analytic
	router.POST("/users/analytic_v1", cAnalytic.PostUserAnalyticV1)
	router.OPTIONS("/users/analytic_v1", baseController.SkipCors)

	// clip
	router.GET("/clips", cClips.GetClips)
	router.POST("/clips", cClips.CreateClip)
	router.OPTIONS("/clips", baseController.SkipCors)

	// animepahe
	router.GET("/animes_map", cAnimes.GetAnimesMap)
	router.GET("/search_anime", cAnimes.GetSearchAnime)

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

	router.Run(":" + port)
}
