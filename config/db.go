package config

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

var DB *mgo.Database
var Consumptions *mgo.Collection
var Fuels *mgo.Collection
var Scrapers *mgo.Collection
var Watcher *mgo.Collection
var Covid *mgo.Collection

func init() {
	s, err := mgo.Dial("mongodb://localhost:27017/hassio")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("hassio")
	Consumptions = DB.C("consumptions")
	Scrapers = DB.C("scrapers")
	Watcher = DB.C("watcher")
	Covid = DB.C("covid")

	fmt.Println("You connected to your mongo database.")
}
