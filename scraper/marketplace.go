package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"../config"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
)

type Watcher struct {
	Item  string // `json:"item" bson:"item"`
	Link  string // `json:"link" bson:"link"`
	Price string // `json:"price" bson:"price"`
}

func GetJofogas(item string) {
	url := fmt.Sprintf("https://www.jofogas.hu/magyarorszag?f=p&q=%s", item)

	scrapeItem(url, ".general-item", ".subject", ".price-value")
}

func GetHvapro(item string) {
	url := fmt.Sprintf("https://hardverapro.hu/aprok/keres.php?stext=%s", item)

	scrapeItem(url, ".media", ".uad-title a", ".uad-price")
}

func checkWatcherItem(link string) bool {
	item := Watcher{}
	err := config.Watcher.Find(bson.M{"link": link}).One(&item)
	if err != nil {
		return false
	}

	return true
}

func insertWatcherItem(title string, link string, price string) error {
	err := config.Watcher.Insert(bson.M{"title": title, "price": price, "link": link})
	if err != nil {
		return err
	}

	return nil
}

func scrapeItem(url string, itemQuery string, titleQuery string, priceQuery string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	doc.Find(itemQuery).Each(func(_ int, s *goquery.Selection) {
		title := s.Find(titleQuery).First()
		name := strings.TrimSpace(title.Text())
		link, _ := title.Attr("href")
		price := strings.TrimSpace(s.Find(priceQuery).First().Text())
		isThere := checkWatcherItem(link)
		if !isThere {
			err := insertWatcherItem(name, link, price)
			log.Printf("Item: %s - %s\n", name, price)
			// TODO: Turn on when actually using
			// slackNotif := fmt.Sprintf("*%s - %s*\n%s", name, price, link)
			// slack.NotifySlack("SLACK_PRESENCE", slackNotif)
			if err != nil {
				log.Printf("Something happened while inserting %s", name)
			}
		}
	})
}
