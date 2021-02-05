package scraper

import (
	"../config"
	"gopkg.in/mgo.v2/bson"
)

type ScraperModel struct {
	Name  string // `json:"name" bson:"name"`
	Value string // `json:"value" bson:"value"`
}

func GetScraperData(name string) (string, error) {
	data := ScraperModel{}

	err := config.Scrapers.Find(bson.M{"name": name}).One(&data)
	if err != nil {
		return name, err
	}
	return data.Value, nil
}

func SaveScraperData(name string, value string) error {
	data := ScraperModel{name, value}
	_, err := config.Scrapers.Upsert(bson.M{"name": name}, &data)
	if err != nil {
		return err
	}
	return nil
}
