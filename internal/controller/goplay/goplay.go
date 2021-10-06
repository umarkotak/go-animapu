package goplay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetProxy(c *gin.Context) {
	goplayPath := c.Request.URL.String()
	goplayFullPath := strings.Replace(goplayPath, "/goplay_api", "https://goplay.co.id/api", -1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", goplayFullPath, nil)
	if err != nil {
		http_req.RenderResponse(c, 422, fmt.Sprintf("error: %v", err))
	}
	req.Header.Add("Authorization", c.Request.Header["Authorization"][0])

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("client.Do: %v", err)
		http_req.RenderResponse(c, 422, fmt.Sprintf("client Do error: %v", err))
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http_req.RenderResponse(c, 422, fmt.Sprintf("ioutil ReadAll error: %v", err))
		return
	}

	var responseBody interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		http_req.RenderResponse(c, 422, fmt.Sprintf("json Unmarshal error: %v", err))
		return
	}

	http_req.RenderResponse(c, 200, responseBody)
}
