package manga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

// GetManga get list of all manga in DB
func GetManga(c *gin.Context) {
	jsonFile, err := os.Open("data/mangas.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	c.JSON(200, result)
}
