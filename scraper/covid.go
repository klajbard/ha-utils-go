package scraper

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/slack"
	"github.com/klajbard/ha-utils-go/utils"
)

type Covid struct {
	Infected int // `json:"infected" bson:"infected"`
	Dead     int // `json:"dead" bson:"dead"`
	Cured    int // `json:"cured" bson:"cured"`
}

// Check hungarian govs official covid page for
// latest update on number of cases and updates
// via Slack if new data is posted
func UpdateCovid() {
	log.Println("[COVID] Querying...")
	covid := &Covid{}
	resp, err := http.Get("https://koronavirus.gov.hu")
	if err != nil {
		utils.NotifyError(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	infectedPest := getNum(doc.Find("#api-fertozott-pest").Text())
	infectedVidek := getNum(doc.Find("#api-fertozott-videk").Text())
	deadPest := getNum(doc.Find("#api-elhunyt-pest").Text())
	deadVidek := getNum(doc.Find("#api-elhunyt-videk").Text())
	curedPest := getNum(doc.Find("#api-gyogyult-pest").Text())
	curedVidek := getNum(doc.Find("#api-gyogyult-videk").Text())

	covid.Infected = infectedPest + infectedVidek
	covid.Dead = deadPest + deadVidek
	covid.Cured = curedPest + curedVidek

	recentCovid, err := getLastCovid()
	if err != nil {
		utils.NotifyError(err)
	}
	delta := sum(covid) - sum(&recentCovid)

	log.Printf("[COVID] Changes: %d\n", delta)

	if !(reflect.DeepEqual(&recentCovid, covid)) {
		result := fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*\n:chart_with_upwards_trend: *%d*", covid.Infected, covid.Dead, covid.Cured, delta)
		slack.NotifySlack("general", result, ":female-doctor:")
		log.Println("[COVID] Updated! New cases " + string(delta))
		insertCovid(covid)
	}
}

func sum(covid *Covid) int {
	return covid.Infected + covid.Cured + covid.Dead
}

func getNum(input string) int {
	trimmed := strings.ReplaceAll(input, " ", "")
	num, err := strconv.Atoi(trimmed)
	if err != nil {
		utils.NotifyError(err)
	}

	return num
}

func insertCovid(covid *Covid) (err error) {
	return config.Covid.Insert(covid)
}

func getLastCovid() (Covid, error) {
	covid := Covid{}

	err := config.Covid.Find(nil).Sort("-_id").Limit(1).One(&covid)
	if err != nil {
		return covid, err
	}

	return covid, nil
}
