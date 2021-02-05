package consumptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"../config"
	"gopkg.in/mgo.v2/bson"
)

type Consumption struct {
	Device string  // `json:"device" bson:"device"`
	Date   string  // `json:"date" bson:"date"`
	Watt   float64 // `json:"watt" bson:"watt"`
}

func AllCons(r *http.Request, d string) ([]Consumption, error) {
	cons := []Consumption{}

	err := config.Consumptions.Find(bson.M{"device": d}).All(&cons)
	if err != nil {
		return nil, err
	}
	return cons, nil
}

func OneCons(r *http.Request, d string, date string) Consumption {
	cons := Consumption{}

	err := config.Consumptions.Find(bson.M{"device": d, "date": date}).One(&cons)
	if err != nil {
		return Consumption{
			Device: d,
			Date:   date,
			Watt:   0.0,
		}
	}
	return cons
}

func UpdateCons(r *http.Request) (Consumption, error) {
	cons := Consumption{}
	json.NewDecoder(r.Body).Decode(&cons)
	fmt.Println(cons)

	if cons.Device == "" || cons.Date == "" || cons.Watt == 0 {
		return cons, errors.New("400. Bad request. All fields must be complete.")
	}

	_, err := config.Consumptions.Upsert(bson.M{"device": cons.Device, "date": cons.Date}, &cons)
	if err != nil {
		return cons, err
	}
	return cons, nil
}
