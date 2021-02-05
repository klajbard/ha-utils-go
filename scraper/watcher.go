package scraper

import (
	"../config"
	"gopkg.in/mgo.v2/bson"
)

type Watcher struct {
	Item  string // `json:"item" bson:"item"`
	Link  string // `json:"link" bson:"link"`
	Price string // `json:"price" bson:"price"`
}

func CheckWatcherItem(link string) bool {
	item := Watcher{}
	err := config.Watcher.Find(bson.M{"link": link}).One(&item)
	if err != nil {
		return false
	}

	return true
}

func InsertWatcherItem(title string, link string, price string) error {
	err := config.Watcher.Insert(bson.M{"title": title, "price": price, "link": link})
	if err != nil {
		return err
	}

	return nil
}
