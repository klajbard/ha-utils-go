package sg

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/julienschmidt/httprouter"
)

func QueryEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, _ := http.NewRequest("GET", "https://www.steamgifts.com", nil)

	req.Header.Set("Referer", "www.steamgifts.com")
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: os.Getenv("SG_SESSID")})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	ga_url, _ := doc.Find(".page__heading + div .giveaway__heading__name").First().Attr("href")
	xsrf_token, _ := doc.Find("input[name=\"xsrf_token\"]").First().Attr("value")

	enterGiveAway(ga_url, xsrf_token)
}

func enterGiveAway(link string, token string) {
	id := strings.Split(link, "/")[2]
	payload := url.Values{
		"xsrf_token": {token},
		"do":         {"entry_insert"},
		"code":       {id},
	}
	req, _ := http.NewRequest("POST", "https://www.steamgifts.com/ajax.php", strings.NewReader(payload.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload.Encode())))
	req.Header.Set("Referer", "www.steamgifts.com"+link)
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: os.Getenv("SG_SESSID")})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	ret, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(ret))
}
