package scraper

import (
	"errors"

	"../config"
	"gopkg.in/mgo.v2/bson"
)

type Fuel struct {
	Text string // `json:"text" bson:"text"`
}

func GetLast() (string, error) {
	fuel := Fuel{}

	err := config.Fuels.Find(nil).Sort("-_id").Limit(1).One(&fuel)
	if err != nil {
		return "", err
	}

	return fuel.Text, nil
}

func SaveFuel(text string) (err error) {
	err = nil
	if text == "" {
		err = errors.New("Text should have length.")
		return err
	}

	err = config.Fuels.Insert(bson.M{"text": text})
	return err
}
