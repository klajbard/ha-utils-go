package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/klajbard/ha-utils-go/slack"
)

func SetState(sensor string, payload interface{}) {
	link := os.Getenv("HASS_URL") + "/api/states/" + sensor
	token := "Bearer " + os.Getenv("HASS_TOKEN")

	payloadString, err := json.Marshal(payload)
	if err != nil {
		NotifyError(err)
	}

	req, err := http.NewRequest("POST", link, strings.NewReader(string(payloadString)))
	if err != nil {
		NotifyError(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(payloadString)))
	req.Header.Set("Authorization", token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		NotifyError(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)
}

func ScrapeFirst(url, query string) string {
	resp, err := http.Get(url)
	if err != nil {
		NotifyError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		NotifyError(err)
	}

	return doc.Find(query).First().Text()
}

func NotifyError(err error) {
	slack.NotifySlack("hass", err.Error(), ":exclamation:")
	PrintError(err)
}

func PrintError(err error) {
	log.Printf("[**Error**] %s", err)
}
