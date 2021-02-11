package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"../utils"
)

// Slack notification uses hooks to send messages
// *group* - The hash for the channel
// Can be extracted from: 'https://hooks.slack.com/services/*group*'
// *text* - The sent text formatted with markup
func NotifySlack(group string, text string) {
	channel := os.Getenv(group)
	if channel == "" {
		utils.PrintError(errors.New("Unable to get ENV variable"))
	}
	reqBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		utils.PrintError(err)
	}

	resp, err := http.Post("https://hooks.slack.com/services/"+channel,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		utils.PrintError(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError(err)
	}

	log.Println(string(body))
}
