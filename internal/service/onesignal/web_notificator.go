package onesignal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// SendWebNotification to send web notification via onesignal
// https://documentation.onesignal.com/reference/create-notification#example-code---create-notification
func SendWebNotification(heading, content string) {
	client := &http.Client{}
	jsonString := fmt.Sprintf(
		`
			{
				"app_id": "%v",
				"contents": {"en": "%v"},
				"headings": {"en": "%v"},
				"included_segments": ["All"]
			}
		`,
		os.Getenv("ONESIGNAL_APP_ID"),
		content,
		heading,
	)
	postBody := []byte(jsonString)
	reqBody := bytes.NewBuffer(postBody)
	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", reqBody)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Basic YzNlZjMyMjQtNmJmNy00ODVlLTliMzItMDdjNjlmMzQyZWNm")
	authVal := fmt.Sprintf("Basic %v", os.Getenv("ONESIGNAL_KEY"))
	req.Header.Set("Authorization", authVal)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

// jsonString := fmt.Sprintf(
// 	`
// 		{
// 			"app_id": "%v",
// 			"contents": {"en": "%v"},
// 			"headings": {"en": "%v"},
// 			"included_segments": ["All"]
// 		}
// 	`,
// 	"46a332f6-fe87-4e33-9b35-636f329e6f6a",
// 	"one piece, enen no shoubutai, boruto, william moriarty, black clover",
// 	"New chapter update!",
// )

// curl --include \
//      --request POST \
//      --header "Content-Type: application/json; charset=utf-8" \
//      --header "Authorization: Basic YzNlZjMyMjQtNmJmNy00ODVlLTliMzItMDdjNjlmMzQyZWNm" \
//      --data-binary "{\"app_id\": \"46a332f6-fe87-4e33-9b35-636f329e6f6a\",
// \"contents\": {\"en\": \"one piece, enen no shoubutai, boruto, william moriarty, black clover\"},
// \"headings\": {\"en\": \"New chapter update\!\"},
// \"included_segments\": [\"All\"]}" \
//      https://onesignal.com/api/v1/notifications
