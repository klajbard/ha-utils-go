package bumpha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/utils"
)

// Update sends a bump for the selected item on "hardverapro"
// *fid* - payload value of 'fidentifier' when bumping
// To gather the fidentifier value, the easiest way is to capture
// from Chrome Dev Tools -> Network tab. Check the 'Preserve log'
// and do the manual bump. When the bump is done search for the
// post request of "felhoz.php", then check the request body
// *name* - name of the item
// The easiest way is to simply copy from the URL
func Update(hvaID string, item config.HaItem) (stop bool) {
	if !shouldBump(item.Start) {
		return
	}
	link := fmt.Sprintf("https://hardverapro.hu/apro/%s/hsz_1-50.html", item.Name)

	resp, err := http.Get(link)
	if err != nil {
		utils.NotifyError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.NotifyError(err)
	}

	lastBump := doc.Find("[title=\"Utolsó UP dátuma\"]").Text()
	pidLink, _ := doc.Find(".row.uad-actions a.btn.btn-secondary.btn-sm").Attr("href")

	re := regexp.MustCompile(`uadid=\d+`)
	pidArr := strings.Split(re.FindString(pidLink), "=")
	if len(pidArr) < 2 {
		return
	}
	pid := pidArr[1]

	if strings.Contains(lastBump, "napja") || lastBumpHours(lastBump) >= item.Hour {
		stop = bumpItem(hvaID, item.Id, pid)
	}
	return
}

func bumpItem(hvaID, fid, pid string) (stop bool) {
	payload := "fidentifier=" + fid
	link := "https://hardverapro.hu/muvelet/apro/felhoz.php?id=" + pid
	req, err := http.NewRequest("POST", link, strings.NewReader(payload))
	if err != nil {
		utils.NotifyError(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload)))
	req.AddCookie(&http.Cookie{Name: "identifier", Value: hvaID})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		utils.NotifyError(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.NotifyError(err)
		return
	}

	if strings.Contains(string(body), "Nincs t\\u00f6bb ingyenes upod m\\u00e1ra!") {
		stop = true
	}
	log.Println("[BUMPHA] " + string(body))
	return
}

func shouldBump(start int) bool {
	now := time.Now().Hour()
	return now >= start
}

func lastBumpHours(lastBump string) int {
	if strings.Contains(lastBump, "órája") {
		re := regexp.MustCompile(`\d`)
		hoursAgo, _ := strconv.Atoi(re.FindString(lastBump))
		return hoursAgo
	}
	return 0
}
