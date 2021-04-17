package config

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

type CounterType struct {
	Name      string `yaml:"name"`
	Available bool   `yaml:"available"`
}

func HandleResetCounter() {
	log.Println("[RESET] Counters reset.")
	err := Counter.Update(bson.M{}, bson.M{"$set": bson.M{"available": true}})
	if err != nil {
		log.Println(err)
	}
}

func GetCounter(name string) *CounterType {
	counter := CounterType{}
	err := Counter.Find(bson.M{"name": name}).One(&counter)
	if err != nil {
		log.Println(err)
	}
	return &counter
}

func UpdateCounter(name string) {
	counter := &CounterType{name, false}
	_, err := Counter.Upsert(bson.M{"name": name}, counter)
	if err != nil {
		log.Println(err)
	}
}
