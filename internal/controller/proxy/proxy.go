package proxy

import (
	"fmt"
	"io"
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
