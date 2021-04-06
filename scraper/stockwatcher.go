package scraper

import (
	"fmt"

	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/slack"
	"github.com/klajbard/ha-utils-go/utils"
	"gopkg.in/mgo.v2/bson"
)

type ScraperModel struct {
	Name  string // `json:"name" bson:"name"`
	Value string // `json:"value" bson:"value"`
}

// Queries for an item availability
// Sends Slack message if it is in stock
func StockWatcher(item *config.ItemStock) {
	result := utils.ScrapeFirst(item.Url, item.Query)
	recentResult, err := getScraperData(item.Name)
	if (err != nil) || recentResult != result {
		slackNotif := fmt.Sprintf("<%s|%s>: %s", item.Url, item.Name, result)
		slack.NotifySlack("general", slackNotif, ":package:")
		err := saveScraperData(item.Name, result)
		if err != nil {
			utils.NotifyError(err)
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

func saveScraperData(name, value string) error {
	data := ScraperModel{name, value}
	_, err := config.Scrapers.Upsert(bson.M{"name": name}, &data)
	return err
}
