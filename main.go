package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./awscost"
	"./bumpha"
	"./dht"
	"./scraper"
	"./sg"
)

func main() {
	tick := time.NewTicker(10 * time.Second)
	tickerCount := 0
	ticker := make(chan int)
	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			i := <-ticker
			scraper.UpdateBestBuy()
			if i%6 == 0 {
				log.Println("Item watcher")
				scraper.GetJofogas("raspberry")
				scraper.GetHvapro("raspberry")
				scraper.GetJofogas("zigbee")
				scraper.GetHvapro("zigbee")
			}
			if i%30 == 0 {
				log.Println("Sg")
				sg.QueryEntry()
			}
			if i%60 == 0 {
				log.Println("DHT")
				dht.ReadDHT(4)
			}
			if i%180 == 0 {
				log.Println("COVID")
				scraper.UpdateCovid()
				log.Println("Bump HVA")
				bumpha.Update(os.Getenv("HVA_ITEM"), "bontatlan_kitvision_escape_hd5w_1080p_akciokamera_3")
			}
			if i%720 == 0 {
				log.Println("Ncore")
				scraper.UpdateNcore()
			}
			if i%1440 == 0 {
				log.Println("Fuel")
				scraper.UpdateFuelPrice()
			}
			if i%4320 == 0 {
				log.Println("Fixer")
				scraper.UpdateCurrencies()
				log.Println("AWS")
				awscost.Update()
			}
			// Temporary turning off
			if false {
				scraper.PicoScraper()
			}
		}
	}()

	for {
		select {
		case <-tick.C:
			ticker <- tickerCount
			tickerCount++
		case <-terminate:
			log.Print("\n\nSIGTERM received. Shutting down...\n")
			return
		}
	}
}
