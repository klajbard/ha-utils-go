package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Slack notification uses hooks to send messages
// *group* - The hash for the channel
// Can be extracted from: 'https://hooks.slack.com/services/*group*'
// *text* - The sent text formatted with markup
func NotifySlack(group string, text string) {
	channel := os.Getenv(group)
	if channel == "" {
		log.Fatalln("Unable to get ENV variable")
	}
	reqBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://hooks.slack.com/services/"+channel,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}
