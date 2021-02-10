package main

import (
	"log"
	"os"
	"time"

	"./awscost"
	"./bumpha"
	"./dht"
	"./scraper"
	"./sg"
)

func main() {
	tick := time.NewTicker(5 * time.Second)
	tickerCount := 0
	ticker := make(chan int)

	go func() {
		for {
			i := <-ticker
			scraper.PicoScraper()
			scraper.UpdateBestBuy()
			if i%12 == 0 {
				log.Println("Item watcher")
				scraper.GetJofogas("raspberry")
				scraper.GetHvapro("raspberry")
			}
			if i%60 == 0 {
				log.Println("Sg")
				sg.QueryEntry()
			}
			if i%120 == 0 {
				log.Println("DHT")
				dht.ReadDHT(4)
			}
			if i%360 == 0 {
				log.Println("COVID")
				scraper.UpdateCovid()
				log.Println("Bump HVA")
				bumpha.Update(os.Getenv("HVA_ITEM"), "bontatlan_kitvision_escape_hd5w_1080p_akciokamera_3")
			}
			if i%1440 == 0 {
				log.Println("Ncore")
				scraper.UpdateNcore()
			}
			if i%2880 == 0 {
				log.Println("Fuel")
				scraper.UpdateFuelPrice()
			}
			if i%8640 == 0 {
				log.Println("Fixer")
				scraper.UpdateCurrencies()
				log.Println("AWS")
				awscost.Update()
			}
		}
	}()

	for {
		select {
		case <-tick.C:
			ticker <- tickerCount
			tickerCount++
		}
	}
}
