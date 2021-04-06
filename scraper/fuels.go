package scraper

import (
	"errors"
	"log"

	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/slack"
	"github.com/klajbard/ha-utils-go/utils"
	"gopkg.in/mgo.v2/bson"
)

type Fuel struct {
	Text string // `json:"text" bson:"text"`
}

// Queries for average fuel prices and updates
// via Slack if theres any new post
func UpdateFuelPrice() {
	result := utils.ScrapeFirst("https://holtankoljak.hu/uzemanyag_arvaltozasok#tartalom", "#Holtankoljak_cikk_leaderboard_top_1 ~ .row .container a")
	recentResult := getRecentFuelPrice()
	if recentResult != result {
		slack.NotifySlack("fuel", result, ":fuelpump:")
		log.Println("[FUEL] New post with price updates is available.")
		err := saveFuelPrice(result)
		if err != nil {
			utils.NotifyError(err)
		}
	}
}

func getRecentFuelPrice() string {
	fuel := Fuel{}

	err := config.Fuels.Find(nil).Sort("-_id").Limit(1).One(&fuel)
	if err != nil {
		utils.NotifyError(err)
	}

	return fuel.Text
}

func saveFuelPrice(text string) (err error) {
	if text == "" {
		return errors.New("Text should have length")
	}

	return config.Fuels.Insert(bson.M{"text": text})
}
