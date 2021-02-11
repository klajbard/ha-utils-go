package scraper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"../config"
	"../slack"
	"../utils"
	"github.com/PuerkitoBio/goquery"

	"gopkg.in/mgo.v2/bson"
)

type BestBuy struct {
	Anchor string // `json:"anchor" bson:"anchor"`
	Text   string // `json:"text" bson:"text"`
	Url    string // `json:"url" bson:"url"`
}

// Checks for the "#akcio" hashtag of
// the bestbuy topic from prohardver and updates
// via Slack if new post is available
func UpdateBestBuy() {
	link := "https://prohardver.hu/tema/bestbuy_topik_akcio_ajanlasakor_akcio_hashtag_kote/friss.html"
	resp, err := http.Get(link)
	if err != nil {
		utils.PrintError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.PrintError(err)
	}
	doc.Find(".msg-list:not(.thread-content) .media").Each(func(id int, s *goquery.Selection) {
		anchorElem := s.Find(".msg-head-author .msg-num a")
		anchor := anchorElem.Text()
		url, _ := anchorElem.Attr("href")
		msg := s.Find(".msg-content").Text()
		re := regexp.MustCompile(`(?i)#akci[oó]`)
		match := re.Match([]byte(msg))

		if match || s.HasClass("msg-featured") {
			foundBB := findBestBuy(anchor)
			if anchor != foundBB.Anchor {
				log.Println("[BESTBUY] New post")
				notif := fmt.Sprintf("<https://prohardver.hu%s|#akcio> *%s* \n%s", url, anchor, msg)
				slack.NotifySlack("SLACK_BESTBUY", notif)
				insertBestBuy(BestBuy{anchor, msg, url})
			}
		}
	})
}

func findBestBuy(anchor string) (bb BestBuy) {
	err := config.BestBuy.Find(bson.M{"anchor": anchor}).One(&bb)
	if err != nil && err.Error() != "not found" {
		utils.PrintError(err)
	}
	return
}

func insertBestBuy(bb BestBuy) {
	err := config.BestBuy.Insert(bb)
	if err != nil {
		utils.PrintError(err)
	}
}
