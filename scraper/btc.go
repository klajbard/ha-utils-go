package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/klajbard/ha-utils-go/types"
	"github.com/klajbard/ha-utils-go/utils"
)

func UpdateBTC() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=BTC", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("BTC_API"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	btcData := types.BTCQuery{}
	err = json.Unmarshal([]byte(string(respBody)), &btcData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(btcData.Data.Btc.Quote.Usd.Price)
	payloadBtc := utils.Sensor{
		State: fmt.Sprintf("%.2f", btcData.Data.Btc.Quote.Usd.Price),
		Attributes: utils.Attributes{
			Friendly_name:       "BTC price",
			Unit_of_measurement: "$",
			Icon:                "mdi:currency-usd",
		},
	}
	utils.SetState("sensor.btc", payloadBtc)
}
