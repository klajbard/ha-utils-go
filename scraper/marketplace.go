package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"../config"
	"../slack"
	"../utils"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
)

type Watcher struct {
	Item  string // `json:"item" bson:"item"`
	Link  string // `json:"link" bson:"link"`
	Price string // `json:"price" bson:"price"`
}

// Queries Jofogas for a specific 'item' and
// sends message to Slack if theres new item
func GetJofogas(item string) {
	url := fmt.Sprintf("https://www.jofogas.hu/magyarorszag?f=p&q=%s", item)

	scrapeItem(url, ".general-item", ".subject", ".price-value")
}

// Queries Hardverapro for a specific 'item' and
// sends message to Slack if theres new item
func GetHvapro(item string) {
	url := fmt.Sprintf("https://hardverapro.hu/aprok/keres.php?stext=%s", item)

	scrapeItem(url, ".media", ".uad-title a", ".uad-price")
}

func checkWatcherItem(link string) bool {
	item := Watcher{}
	err := config.Watcher.Find(bson.M{"link": link}).One(&item)
	return err == nil
}

func insertWatcherItem(title, link, price string) error {
	return config.Watcher.Insert(bson.M{"title": title, "price": price, "link": link})
}

func scrapeItem(url, itemQuery, titleQuery, priceQuery string) {
	resp, err := http.Get(url)
	if err != nil {
		utils.PrintError(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.PrintError(err)
		return
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
			slackNotif := fmt.Sprintf("<%s|*%s - %s*>\n", link, name, price)
			slack.NotifySlack("SLACK_PRESENCE", slackNotif)
			if err != nil {
				log.Printf("Something happened while inserting %s", name)
			}
		}
	})
}
