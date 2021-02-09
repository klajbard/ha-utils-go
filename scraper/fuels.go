package scraper

import (
	"errors"
	"log"

	"../config"
	"../slack"
	"../utils"
	"gopkg.in/mgo.v2/bson"
)

type Fuel struct {
	Text string // `json:"text" bson:"text"`
}

func UpdateFuelPrice() {
	result := utils.ScrapeFirst("https://holtankoljak.hu/uzemanyag_arvaltozasok#tartalom", "#Holtankoljak_cikk_leaderboard_top_1 ~ .row .container a")
	recentResult := getRecentFuelPrice()
	if recentResult != result {
		slack.NotifySlack("SLACK_SCRAPER", result)
		log.Println("[FUEL] New post with price updates is available.")
		err := saveFuelPrice(result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func getRecentFuelPrice() string {
	fuel := Fuel{}

	err := config.Fuels.Find(nil).Sort("-_id").Limit(1).One(&fuel)
	if err != nil {
		log.Fatalln(err)
	}

	return fuel.Text
}

func saveFuelPrice(text string) (err error) {
	err = nil
	if text == "" {
		err = errors.New("Text should have length")
		return err
	}

	err = config.Fuels.Insert(bson.M{"text": text})
	return err
}
