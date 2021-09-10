package mangadex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/lib/utils"
)

func GetProxy(c *gin.Context) {
	mangadexPath := c.Request.URL.String()
	mangadexFullpath := strings.Replace(mangadexPath, "/mangadex", "https://api.mangadex.org", -1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", mangadexFullpath, nil)
	if err != nil {
		utils.RenderResponse(c, 422, fmt.Sprintf("error: %v", err))
	}
	req.Header.Add("Authorization", c.Request.Header["Authorization"][0])

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var responseBody interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		utils.RenderResponse(c, 422, fmt.Sprintf("error: %v", err))
	}

	utils.RenderResponse(c, 200, responseBody)
}
