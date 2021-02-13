// Bumps an item on hardverapro
package bumpha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../utils"
	"github.com/PuerkitoBio/goquery"
)

// Sends a bump for the selected item on "hardverapro"
// *fid* - payload value of 'fidentifier' when bumping
// To gather the fidentifier value, the easiest way is to capture
// from Chrome Dev Tools -> Network tab. Check the 'Preserve log'
// and do the manual bump. When the bump is done search for the
// post request of "felhoz.php", then check the request body
// *name* - name of the item
// The easiest way is to simply copy from the URL
func Update(fid, name string) {
	link := fmt.Sprintf("https://hardverapro.hu/apro/%s/hsz_1-50.html", name)

	resp, err := http.Get(link)
	if err != nil {
		utils.PrintError(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utils.PrintError(err)
	}

	last_bump := doc.Find("[title=\"Utolsó UP dátuma\"]").Text()
	pid_link, _ := doc.Find(".row.uad-actions a.btn.btn-secondary.btn-sm").Attr("href")

	re := regexp.MustCompile(`uadid=\d+`)
	pid := strings.Split(re.FindString(pid_link), "=")[1]

	if strings.Contains(last_bump, "napja") {
		bumpItem(fid, pid)
	} else if strings.Contains(last_bump, "órája") {
		re := regexp.MustCompile(`\d`)
		hours_ago, _ := strconv.Atoi(re.FindString(last_bump))
		if hours_ago > 2 {
			bumpItem(fid, pid)
		}
	}
}

func bumpItem(fid, pid string) {
	payload := "fidentifier=" + fid
	link := "https://hardverapro.hu/muvelet/apro/felhoz.php?id=" + pid
	req, err := http.NewRequest("POST", link, strings.NewReader(payload))
	if err != nil {
		utils.PrintError(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload)))
	req.AddCookie(&http.Cookie{Name: "identifier", Value: os.Getenv("HVA_ID")})

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		utils.PrintError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError(err)
	}

	log.Println("[BUMPHA] " + string(body))
}
