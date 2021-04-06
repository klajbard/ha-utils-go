package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/slack"
	"github.com/klajbard/ha-utils-go/utils"
	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Price int    `yaml:"price"`
}

func QueryArukereso(url string) {
	var minPrice int

	resp, err := http.Get(url)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	title := strings.TrimSpace(doc.Find("h1.hidden-xs").Text())
	item := getItem(title)

	doc.Find(".optoffer.device-desktop").Each(func(_ int, s *goquery.Selection) {
		price, _ := s.Find("[itemprop=\"price\"]").Attr("content")
		priceInt, _ := strconv.Atoi(price)
		if minPrice == 0 || priceInt < minPrice {
			minPrice = priceInt
		}
	})
	if minPrice < item.Price || item.Price == 0 {
		saveItem(title, url, minPrice)
		text := fmt.Sprintf("<%s|*%s - %d*>\n", url, title, minPrice)
		slack.NotifySlack("general", text, ":desktop_computer:")
	}
}

func getItem(name string) (item Item) {
	err := config.Arukereso.Find(bson.M{"name": name}).Limit(1).One(&item)
	if err != nil {
		log.Println(err)
	}

	return item
}

func saveItem(name, url string, price int) {
	item := &Item{name, url, price}
	_, err := config.Arukereso.Upsert(bson.M{"name": name}, item)
	if err != nil {
		log.Println(err)
	}
}
