package apple

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

type (
	AppleCallback struct {
		Code    string
		IdToken string
		User    AppleUser
	}

	AppleUser struct {
		ID    string
		Name  AppleName
		Email string
	}

	AppleName struct {
		FirstName string
		LastName  string
	}
)

func Callback(c *gin.Context) {
	appleUser := AppleUser{}
	if c.PostForm("user") != "" {
		json.Unmarshal([]byte(c.PostForm("user")), &appleUser)
	}

	appleCallback := AppleCallback{
		Code:    c.PostForm("code"),
		IdToken: c.PostForm("id_token"),
		User:    appleUser,
	}

	http_req.RenderResponse(c, 200, appleCallback)
}

func CallbackRedirect(c *gin.Context) {
	appleUser := AppleUser{}
	if c.PostForm("user") != "" {
		json.Unmarshal([]byte(c.PostForm("user")), &appleUser)
	}

	appleCallback := AppleCallback{
		Code:    c.PostForm("code"),
		IdToken: c.PostForm("id_token"),
		User:    appleUser,
	}

	baseUrl := "https://animapu.netlify.app/goplay/apple/account/return"
	urlPath := "goplay/apple/account/return"
	urlParams := fmt.Sprintf("code=%v&id_token=%v", appleCallback.Code, appleCallback.IdToken)
	finalUrl := fmt.Sprintf("%v/%v?%v", baseUrl, urlPath, urlParams)

	c.Redirect(http.StatusMovedPermanently, finalUrl)
}
