package scraper

import (
	"fmt"
	"log"

	"../config"
	"../slack"
	"../utils"
	"gopkg.in/mgo.v2/bson"
)

type ScraperModel struct {
	Name  string // `json:"name" bson:"name"`
	Value string // `json:"value" bson:"value"`
}

func PicoScraper() {
	url := "https://www.optimusdigital.ro/en/raspberry-pi-boards/12024-raspberry-pi-pico-728886755172.html"
	result := utils.ScrapeFirst(url, "#quantityAvailable")
	recentResult, err := getScraperData("pico")
	if (err != nil) || recentResult != result {
		slackNotif := fmt.Sprintf("Presence checker changed on Raspberry Pico: %s \n%s", result, url)
		slack.NotifySlack("SLACK_PRESENCE", slackNotif)
		err := saveScraperData("pico", result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func getScraperData(name string) (string, error) {
	data := ScraperModel{}

	err := config.Scrapers.Find(bson.M{"name": name}).One(&data)
	if err != nil {
		return name, err
	}
	return data.Value, nil
}

func saveScraperData(name string, value string) error {
	data := ScraperModel{name, value}
	_, err := config.Scrapers.Upsert(bson.M{"name": name}, &data)
	if err != nil {
		return err
	}
	return nil
}
