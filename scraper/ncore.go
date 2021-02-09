package scraper

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"

	"../utils"
	"github.com/PuerkitoBio/goquery"
)

var ncoreMap = []string{
	"ncore_daily_rank",
	"ncore_weekly_rank",
	"ncore_monthly_rank",
	"ncore_last_month_rank",
	"ncore_allowed",
	"ncore_active",
	"ncore_hitnrun_possible",
	"ncore_hitnrun_month",
	"ncore_possible_hit_n_run",
}

func UpdateNcore() {
	sessid := ""
	login := fmt.Sprintf("nev=%s&pass=%s", os.Getenv("NCORE_USERNAME"), os.Getenv("NCORE_PASSWORD"))
	jar, _ := cookiejar.New(nil)
	req, err := http.NewRequest("POST", "https://ncore.pro/login.php", strings.NewReader(login))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	req.Header.Set("Origin", "https://ncore.pro")
	req.Header.Set("Referer", "https://ncore.pro/login.php")
	req.Header.Add("Content-Length", strconv.Itoa(len(login)))

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // https://stackoverflow.com/a/38150816
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PHPSESSID" {
			sessid = cookie.Value
		}
	}

	log.Println("[NCORE] Updating profile statistics")

	getData(sessid)
}

func getData(sessid string) {
	result := map[string]string{}
	req, err := http.NewRequest("POST", "https://ncore.pro/hitnrun.php", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: sessid})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	currentData := ""

	doc.Find(".fobox_tartalom div:first-of-type .dt, .fobox_tartalom div:first-of-type .dd").Each(func(index int, s *goquery.Selection) {
		if index%2 == 0 {
			currentData = strings.TrimSuffix(s.Text(), ":")
		} else {
			result[currentData] = s.Text()
			updateHassio(currentData, s.Text(), ncoreMap[index/2])
		}
	})
}

func updateHassio(name string, value string, sensorName string) {
	sensor := utils.Sensor{
		State: value,
		Attributes: utils.Attributes{
			Friendly_name: name,
			Icon:          "mdi:server-network",
		},
	}
	utils.SetState("sensor."+sensorName, sensor)
}
