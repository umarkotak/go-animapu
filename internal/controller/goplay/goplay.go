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
	goplayFullPath := strings.Replace(goplayPath, "/goplay_api_integration", "https://r7lk8n0srbams5t7sl7q.goplay.co.id/api", -1)
	goplayFullPath = strings.Replace(goplayPath, "/goplay_api", "https://goplay.co.id/api", -1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", goplayFullPath, nil)
	if err != nil {
		http_req.RenderResponse(c, 422, fmt.Sprintf("error: %v", err))
	}
	// req.Header.Add("Authorization", c.Request.Header["Authorization"][0])
	req.Header.Add("Cookie", "_ga=GA1.3.1226319411.1594254951; afUserId=6f0973c5-c1be-494e-8d9a-3838f72b48e6-p; WZRK_G=0eca1152d3af476987c83e31cb8f6b10; _gcl_au=1.1.384206423.1626833327; _fbp=fb.2.1628051317781.1640788811; OptanonAlertBoxClosed=2021-08-18T09:45:03.106Z; _gid=GA1.3.743409497.1633427768; AF_SYNC=1633427769730; gp_daily_session=9f317978-8cad-4ea7-95fb-bfdaae728819; _gorilla_csrf=MTYzMzU3NzgyMnxJakZHVG1jd05sVlJTellyYVVaS2F6UmhkVUZ5Tld0VVpWRjFWRXN4TTBWSWNXdGlkbmxvY0VSbmFuYzlJZ289fL6xw6Ka6qAIYC3RL7rWtAONY7UkgPvM8ohbSRbaX5Eb; WZRK_S_R5Z-568-WK5Z=%%7B%%22p%%22%%3A1%%2C%%22s%%22%%3A1633578609%%2C%%22t%%22%%3A1633578609%%7D; session=UYCCCGDDC2DAJ5DGNSYTUXDRSV55PHYQF3L36UGETIASUZGH3MN6K7HWQGA6SJJXWFFN6Z6LOPJ6VPHPRFKRYSU3KHNHOEGWIW47VAA; _gat_UA-162768158-1=1; OptanonConsent=isGpcEnabled=0&datestamp=Thu+Oct+07+2021+10%%3A54%%3A42+GMT%%2B0700+(Western+Indonesia+Time)&version=6.22.0&geolocation=%%3B&isIABGlobal=false&hosts=&landingPath=NotLandingPage&groups=C0001%%3A1%%2CC0002%%3A0%%2CC0004%%3A0&AwaitingReconsent=false&consentId=6e1f1f30-922a-46e4-a056-2b3ce7bc16c7&interactionCount=0; WZRK_S_K5Z-568-WK5Z=%%7B%%22p%%22%%3A5%%2C%%22s%%22%%3A1633577823%%2C%%22t%%22%%3A1633578882%%7D; gp_fgp=7e6e573e30698121f1b2d5018988697d%%3Aed47c4e57b84711139a2f0ccd7a678f0f5225160942464cc78eadc8d2ade4a20")

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
