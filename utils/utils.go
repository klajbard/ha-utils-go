package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func SetState(sensor string, payload interface{}) {
	link := os.Getenv("HASS_URL") + "/api/states/" + sensor
	token := "Bearer " + os.Getenv("HASS_TOKEN")

	payloadString, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", link, strings.NewReader(string(payloadString)))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(payloadString)))
	req.Header.Set("Authorization", token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
