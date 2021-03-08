package config

import (
	"log"

	"gopkg.in/mgo.v2"
)

var DB *mgo.Database
var Fuels *mgo.Collection
var Scrapers *mgo.Collection
var Watcher *mgo.Collection
var Covid *mgo.Collection
var AWS *mgo.Collection
var BestBuy *mgo.Collection
var Arukereso *mgo.Collection

func init() {
	s, err := mgo.Dial("mongodb://localhost:27017/hassio")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("hassio")
	Scrapers = DB.C("scrapers")
	Fuels = DB.C("fuels")
	Watcher = DB.C("watcher")
	Covid = DB.C("covid")
	AWS = DB.C("aws")
	BestBuy = DB.C("bestbuy")
	Arukereso = DB.C("arukereso")

	log.Println("MongoDB connected")
}
