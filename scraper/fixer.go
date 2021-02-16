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
)

// Queries fixer.io for currencies HUF based and
// updates homeassistant corresponding sensor values
// sensor.usd, sensor.eur, sensor.gbp
func UpdateCurrencies() {
	out := struct {
		Success   bool
		Timestamp int64
		Base      string
		Date      string
		Rates     map[string]float64
	}{}

	base := []string{"EUR", "USD", "GBP"}
	target := "HUF"

	link := fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%s&base=EUR", os.Getenv("FIXERAPI"))
	resp, err := http.Get(link)
	if err != nil {
		utils.PrintError(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError(err)
		return
	}

	err = json.Unmarshal([]byte(string(body)), &out)
	if err != nil {
		utils.PrintError(err)
		return
	}

	rates := out.Rates

	log.Println("[FIXER] Updated currencies")

	for _, val := range base {
		rate := rates[target] / rates[val]
		updateHassioScraper(val, target, rate)
	}
}

func updateHassioScraper(name, target string, value float64) {
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
