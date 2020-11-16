package manga

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	sUser "github.com/umarkotak/go-animapu/internal/service/user"
)

// RegisterUserFirebase get list of all manga in DB
func RegisterUserFirebase(c *gin.Context) {
	type RequestParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var json RequestParams
	c.BindJSON(&json)

	userDataBase := models.UserData{
		Username: json.Username,
		Password: json.Password,
	}

	userData, err := sUser.RegisterService(userDataBase)

	var response gin.H
	var statusCode int
	if err != nil {
		response = gin.H{
			"username": userDataBase.Username,
			"message":  err.Error(),
		}
		statusCode = 400
	} else {
		response = gin.H{
			"username": userData.Username,
			"message":  "register success",
		}
		statusCode = 200
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(statusCode, response)
}

// LoginUser run update manga
func LoginUser(c *gin.Context) {
	type RequestParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var json RequestParams
	c.BindJSON(&json)

	userDataBase := models.UserData{
		Username: json.Username,
		Password: json.Password,
	}

	userData, err := sUser.LoginService(userDataBase)

	var response gin.H
	var statusCode int
	if err != nil {
		response = gin.H{
			"message": err.Error(),
		}
		statusCode = 400
	} else {
		response = gin.H{
			"username":    userData.Username,
			"login_token": userData.LoginToken,
			"message":     "login success",
		}
		statusCode = 200
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(statusCode, response)
}

// GetDetailFirebase get manga from firebase
func GetDetailFirebase(c *gin.Context) {
	auth := c.Request.Header["Authorization"][0]
	userData, err := sUser.DetailService(auth)

	var response gin.H
	var statusCode int
	if err != nil {
		response = gin.H{
			"message": err.Error(),
		}
		statusCode = 400
	} else {
		response = gin.H{
			"username":       userData.Username,
			"read_histories": userData.ReadHistories,
			"message":        "success",
		}
		statusCode = 200
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(statusCode, response)
}

// LogReadHistories update user last read histories
func LogReadHistories(c *gin.Context) {
	auth := c.Request.Header["Authorization"][0]

	type RequestParams struct {
		LastChapter string `json:"last_chapter"`
		MangaTitle  string `json:"manga_title"`
	}
	var json RequestParams
	c.BindJSON(&json)

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	lastChapter, _ := strconv.Atoi(json.LastChapter)
	readHistory := models.ReadHistory{
		MangaTitle:    json.MangaTitle,
		LastChapter:   lastChapter,
		LastReadTime:  now.Format(time.RFC3339),
		LastReadTimeI: now.Unix(),
	}

	userData, err := sUser.DetailService(auth)

	if userData.ReadHistories == nil {
		userData.ReadHistories = make(map[string]*models.ReadHistory)
	}
	userData.ReadHistories[readHistory.MangaTitle] = &readHistory

	fmt.Println(userData, readHistory)

	userData, err = sUser.RecordLastReadHistory(userData)

	var response gin.H
	var statusCode int
	if err != nil {
		response = gin.H{
			"message": err.Error(),
		}
		statusCode = 400
	} else {
		response = gin.H{
			"username":       userData.Username,
			"read_histories": userData.ReadHistories,
			"message":        "success",
		}
		statusCode = 200
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(statusCode, response)
}

// SkipCors skip cors
func SkipCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, "CORS OK")
}
