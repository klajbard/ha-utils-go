package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SetState(sensor string, payload interface{}) {
	link := os.Getenv("HASS_URL") + "/api/states/" + sensor
	token := "Bearer " + os.Getenv("HASS_TOKEN")

	payloadString, err := json.Marshal(payload)
	if err != nil {
		PrintError(err)
	}

	req, err := http.NewRequest("POST", link, strings.NewReader(string(payloadString)))
	if err != nil {
		PrintError(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(payloadString)))
	req.Header.Set("Authorization", token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		PrintError(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

func ScrapeFirst(url, query string) string {
	resp, err := http.Get(url)
	if err != nil {
		PrintError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		PrintError(err)
	}

	return doc.Find(query).First().Text()
}

func PrintError(err error) {
	fmt.Printf("[**Error**] %s", err)
}
