package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetProxy(c *gin.Context) {
	currPath := c.Request.URL.String()
	splitPath := strings.Split(currPath, "/proxy/")
	if len(splitPath) != 2 {
		http_req.RenderResponse(c, 422, "invalid path format")
		return
	}

	targetUrl := splitPath[1]
	logrus.Infoln("URL:", targetUrl)

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		http_req.RenderResponse(c, 422, fmt.Sprintf("error: %v", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("client.Do: %v", err)
		http_req.RenderResponse(c, 422, fmt.Sprintf("client Do error: %v", err))
		return
	}

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	http_req.RenderResponse(c, 422, fmt.Sprintf("ioutil ReadAll error: %v", err))
	// 	return
	// }

	copyHeader(c.Writer.Header(), resp.Header)
	io.Copy(c.Writer, resp.Body)
	http_req.RenderResponse(c, 200, nil)
}

func ImageProxy(c *gin.Context) {
	currPath := c.Request.URL.String()
	splitPath := strings.Split(currPath, "/image_proxy/")
	if len(splitPath) != 2 {
		http_req.RenderResponse(c, 422, "invalid path format")
		return
	}

	targetUrl := splitPath[1]

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"98\", \"Google Chrome\";v=\"98\"")
	req.Header.Set("Referer", "https://m.mangabat.com/")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	c.Writer.Write(body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			if k == "X-Frame-Options" || k == "x-frame-options" {
				continue
			}
			dst.Add(k, v)
		}
	}
}
