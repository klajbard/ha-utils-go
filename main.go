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
	tick := time.NewTicker(3 * time.Second)
	tickerCount := 0
	ticker := make(chan int)
	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer handleException()
		for {
			i := <-ticker
			go scraper.UpdateBestBuy()
			go stockWatcher()
			if i%6 == 0 {
				log.Println("Item watcher")
				go handleMarketplace()
			}
			if i%30 == 0 {
				if os.Getenv("SG_SESSID") != "" {
					log.Println("Sg")
					go sg.QueryEntry()
				}
			}
			if i%60 == 0 {
				log.Println("DHT")
				go dht.ReadDHT(4)
				log.Println("Arukereso")
				go queryArukereso()
			}
			if i%180 == 0 {
				log.Println("COVID")
				scraper.UpdateCovid()
				log.Println("Bump HVA")
				go handleHABump()
			}
			if i%720 == 0 {
				if os.Getenv("NCORE_USERNAME") != "" && os.Getenv("NCORE_PASSWORD") != "" {
					log.Println("Ncore")
					go scraper.UpdateNcore()
				}
			}
			if i%1440 == 0 {
				log.Println("Fuel")
				go scraper.UpdateFuelPrice()
			}
			if i%4320 == 0 {
				if os.Getenv("FIXERAPI") != "" {
					log.Println("Fixer")
					go scraper.UpdateCurrencies()
				}
				go awscost.Update()
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
	for _, user := range config.Conf.HaBump {
		for _, item := range user.Items {
			bumpha.Update(user.Identifier, item.Id, item.Name)
		}
	}
}

func stockWatcher() {
	for _, item := range config.Conf.StockWatcher {
		scraper.StockWatcher(&item)
	}
}

func queryArukereso() {
	for _, item := range config.Conf.Arukereso {
		scraper.QueryArukereso(item.Url)
	}
}
