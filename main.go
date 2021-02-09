package main

import (
	"./awscost"
	"./bumpha"
	"./consumptions"
	"./scraper"
	"./sg"
)

func main() {
	bumpha.Update("fidentifier123123123123", "item1")
	if false {
		awscost.Update()
		consumptions.GetAllCons("device-name")
		consumptions.GetCons("device-name", "2021-02-08")
		consumptions.PutCons("device-name", "2021-02-08", 80.18)
		scraper.UpdateCovid()
		scraper.UpdateCurrencies()
		scraper.UpdateFuelPrice()
		scraper.UpdateNcore()
		scraper.PicoScraper()
		scraper.GetJofogas("samsung")
		scraper.GetHvapro("samsung")
		sg.QueryEntry()
	}
}
