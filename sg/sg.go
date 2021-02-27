package sg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"../utils"
	"github.com/PuerkitoBio/goquery"
)

// Enters the upcoming SG giveaway if there is
// enought points left.
func QueryEntry() {
	req, err := http.NewRequest("GET", "https://www.steamgifts.com", nil)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	req.Header.Set("Referer", "www.steamgifts.com")
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: os.Getenv("SG_SESSID")})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		utils.NotifyError(err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	ga_url, _ := doc.Find(".page__heading + div .giveaway__heading__name").First().Attr("href")
	xsrf_token, _ := doc.Find("input[name=\"xsrf_token\"]").First().Attr("value")

	enterGiveAway(ga_url, xsrf_token)
}

func enterGiveAway(link, token string) {
	id := strings.Split(link, "/")[2]
	payload := url.Values{
		"xsrf_token": {token},
		"do":         {"entry_insert"},
		"code":       {id},
	}
	req, err := http.NewRequest("POST", "https://www.steamgifts.com/ajax.php", strings.NewReader(payload.Encode()))
	if err != nil {
		utils.NotifyError(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload.Encode())))
	req.Header.Set("Referer", "www.steamgifts.com"+link)
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: os.Getenv("SG_SESSID")})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		utils.NotifyError(err)
		return
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.NotifyError(err)
	}
	fmt.Println(string(ret))
}
