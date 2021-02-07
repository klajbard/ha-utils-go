package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../utils"
	"github.com/julienschmidt/httprouter"
)

func ScrapeFixer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	out := struct {
		Success   string
		Timestamp string
		Base      string
		Date      string
		Rates     map[string]float64
	}{}

	base := []string{"EUR", "USD", "GBP"}
	target := "HUF"

	link := fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%s&base=EUR", os.Getenv("FIXERAPI"))
	resp, err := http.Get(link)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	err = json.Unmarshal([]byte(string(body)), &out)
	if err != nil {
		print(err)
	}

	rates := out.Rates

	for _, val := range base {
		rate := rates[target] / rates[val]
		updateHassioScraper(val, rate, target)
	}
}

func updateHassioScraper(name string, value float64, target string) {
	valueString := strconv.FormatFloat(value, 'f', 2, 64)
	nameLowercase := strings.ToLower(name)

	sensor := utils.Sensor{
		State: valueString,
		Attributes: utils.Attributes{
			Friendly_name:       name,
			Unit_of_measurement: target,
			Icon:                "mdi:currency-" + nameLowercase,
		},
	}
	utils.SetState("sensor."+nameLowercase, sensor)
}
