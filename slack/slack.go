package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func NotifySlack(group string, text string) {
	channel := os.Getenv(group)
	if channel == "" {
		log.Fatalln("Unable to get ENV variable")
	}
	reqBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		print(err)
	}

	resp, err := http.Post("https://hooks.slack.com/services/"+channel,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	fmt.Println(string(body))
}
