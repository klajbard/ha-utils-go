package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"../slack"
	"github.com/PuerkitoBio/goquery"
	"github.com/julienschmidt/httprouter"
)

func scrapeFirst(url string, query string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return doc.Find(query).First().Text()
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
		isThere := CheckWatcherItem(link)
		if !isThere {
			err := InsertWatcherItem(name, link, price)
			fmt.Printf("Item: %s - %s\n", name, price)
			// slackNotif := fmt.Sprintf("*%s - %s*\n%s", name, price, link)
			// slack.NotifySlack("SLACK_PRESENCE", slackNotif)
			if err != nil {
				fmt.Printf("Something happened while inserting %s", name)
			}
		}
	})
}

func PicoScraper(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := "https://www.optimusdigital.ro/en/raspberry-pi-boards/12024-raspberry-pi-pico-728886755172.html"
	result := scrapeFirst(url, "#quantityAvailable")
	recentResult, err := GetScraperData("pico")
	if (err != nil) || recentResult != result {
		slackNotif := fmt.Sprintf("Presence checker changed on Raspberry Pico: %s \n%s", result, url)
		slack.NotifySlack("SLACK_PRESENCE", slackNotif)
		err := SaveScraperData("pico", result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func GetFuel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := scrapeFirst("https://holtankoljak.hu/uzemanyag_arvaltozasok#tartalom", "#Holtankoljak_cikk_leaderboard_top_1 ~ .row .container a")
	recentResult, err := GetLast()
	if err != nil {
		log.Fatalln(err)
	}
	if recentResult != result {
		slack.NotifySlack("SLACK_SCRAPER", result)
		err := SaveFuel(result)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func GetJofogas(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	item := ps.ByName("item")
	url := fmt.Sprintf("https://www.jofogas.hu/magyarorszag?f=p&q=%s", item)

	scrapeItem(url, ".general-item", ".subject", ".price-value")
}

func GetHvapro(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	item := ps.ByName("item")
	url := fmt.Sprintf("https://hardverapro.hu/aprok/keres.php?stext=%s", item)

	scrapeItem(url, ".media", ".uad-title a", ".uad-price")
}
