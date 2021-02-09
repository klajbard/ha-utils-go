package consumptions

import (
	"log"

	"../config"
	"gopkg.in/mgo.v2/bson"
)

type Consumption struct {
	Device string  // `json:"device" bson:"device"`
	Date   string  // `json:"date" bson:"date"`
	Watt   float64 // `json:"watt" bson:"watt"`
}

func GetAllCons(device string) {
	cons := []Consumption{}
	err := config.Consumptions.Find(bson.M{"device": device}).All(&cons)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cons)
}

func GetCons(device string, date string) {
	cons := Consumption{}
	err := config.Consumptions.Find(bson.M{"device": device, "date": date}).One(&cons)
	if err != nil {
		cons = Consumption{
			Device: device,
			Date:   date,
			Watt:   0.0,
		}
	}

	log.Println(cons)
}

func PutCons(device string, date string, watt float64) {
	cons := Consumption{device, date, watt}
	_, err := config.Consumptions.Upsert(bson.M{"device": cons.Device, "date": cons.Date}, &cons)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("[CONS] Update succeed for " + device)
}
