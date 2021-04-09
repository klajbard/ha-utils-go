package hass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/slack"
	"github.com/klajbard/ha-utils-go/types"
	"github.com/klajbard/ha-utils-go/utils"
	"gopkg.in/mgo.v2/bson"
)

var sensors = []string{
	"zeller_today_kwh",
	"kelbimbo_today",
	"karfiol_today",
	"articsoka_today",
}

type Consumption struct {
	Device string  // `json:"device" bson:"device"`
	Date   string  // `json:"date" bson:"date"`
	Watt   float64 // `json:"watt" bson:"watt"`
}

func SaveCons() {
	for _, sensor := range sensors {
		sensorData := types.SensorCons{}
		link := os.Getenv("HASS_URL") + "/api/states/sensor." + sensor
		token := "Bearer " + os.Getenv("HASS_TOKEN")

		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			utils.NotifyError(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			utils.NotifyError(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.NotifyError(err)
			return
		}
		err = json.Unmarshal([]byte(string(body)), &sensorData)
		if err != nil {
			log.Fatal(err)
		}
		updateCons(strings.Split(sensor, "_")[0], sensorData.State)
	}
}

func updateCons(device, watt string) {
	wattNum, _ := strconv.ParseFloat(watt, 64)
	date := time.Now().Format("2006.01.02")
	cons := Consumption{device, date, wattNum}
	if device == "" || date == "" || wattNum == 0 {
		return
	}

	_, err := config.Consumptions.Upsert(bson.M{"device": device, "date": date}, &cons)
	if err != nil {
		utils.NotifyError(err)
	}
}

func GetYesterdayCons() {
	var sb strings.Builder
	sum := 0.0
	for _, sensor := range sensors {
		device := strings.Split(sensor, "_")[0]
		date := time.Now().AddDate(0, 0, -1).Format("2006.01.02")
		cons := Consumption{}

		err := config.Consumptions.Find(bson.M{"device": device, "date": date}).One(&cons)
		if err != nil {
			log.Println(err)
		} else {
			sum += cons.Watt
			sb.WriteString(fmt.Sprintf("%s: %.1f W\n", cons.Device, cons.Watt))
		}
	}

	price := 37 * sum / 1000
	sb.WriteString(fmt.Sprintf("Total: *%.1f W %.1f Ft*\n", sum, price))
	slack.NotifySlack("consumption", sb.String(), ":desktop_computer:")
}
