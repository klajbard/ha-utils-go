package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./awscost"
	"./bumpha"
	"./config"
	"./dht"
	"./scraper"
	"./sg"
)

func handleException() {
	if e := recover(); e != nil {
		fmt.Println("Recovering from the error: ", e)
	}
}

func main() {
	tick := time.NewTicker(10 * time.Second)
	tickerCount := 0
	ticker := make(chan int)
	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer handleException()
		for {
			i := <-ticker
			scraper.UpdateBestBuy()
			if i%6 == 0 {
				log.Println("Item watcher")
				go handleMarketplace()
			}
			if i%30 == 0 {
				if os.Getenv("SG_SESSID") != "" {
					log.Println("Sg")
					sg.QueryEntry()
				}
			}
			if i%60 == 0 {
				log.Println("DHT")
				dht.ReadDHT(4)
			}
			if i%180 == 0 {
				log.Println("COVID")
				scraper.UpdateCovid()
				if os.Getenv("HVA_ID") != "" {
					log.Println("Bump HVA")
					go handleHABump()
				}
			}
			if i%720 == 0 {
				if os.Getenv("NCORE_USERNAME") != "" && os.Getenv("NCORE_PASSWORD") != "" {
					log.Println("Ncore")
					scraper.UpdateNcore()
				}
			}
			if i%1440 == 0 {
				log.Println("Fuel")
				scraper.UpdateFuelPrice()
			}
			if i%4320 == 0 {
				if os.Getenv("FIXERAPI") != "" {
					log.Println("Fixer")
					scraper.UpdateCurrencies()
				}
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

func handleMarketplace() {
	log.Println(config.Conf.Marketplace.Enabled)
	if !config.Conf.Marketplace.Enabled {
		return
	}
	for _, item := range config.Conf.Marketplace.Jofogas {
		scraper.GetJofogas(item.Name)
	}
	for _, item := range config.Conf.Marketplace.Hardverapro {
		scraper.GetHvapro(item.Name)
	}
}

func handleHABump() {
	for _, item := range config.Conf.HaBump {
		bumpha.Update(item.Id, item.Name)
	}
}
