package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/klajbard/ha-utils-go/awscost"
	"github.com/klajbard/ha-utils-go/bumpha"
	"github.com/klajbard/ha-utils-go/config"
	"github.com/klajbard/ha-utils-go/dht"
	"github.com/klajbard/ha-utils-go/hass"
	"github.com/klajbard/ha-utils-go/scraper"
	"github.com/klajbard/ha-utils-go/sg"
)

func handleException() {
	if e := recover(); e != nil {
		log.Println("Recovering from the error: ", e)
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
			config.Conf.GetConf()
			i := <-ticker
			handleUpdateBB()
			stockWatcher()
			handleNotifyConsumption()
			if i%6 == 0 {
				handleMarketplace()
			}
			if i%60 == 0 {
				handleSG()
				handleUpdateDHT()
				queryArukereso()
				handleBTC()
				handleSaveConsumption()
			}
			if i%180 == 0 {
				handleUpdateCovid()
				handleHABump()
			}
			if i%720 == 0 {
				handleUpdateNcore()
			}
			if i%1440 == 0 {
				handleUpdateFuel()
			}
			if i%4320 == 0 {
				handleUpdateFixer()
				handleAwsCost()
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

func handleSaveConsumption() {
	if !config.Conf.Enable.LogConsumption || os.Getenv("HASS_URL") == "" || os.Getenv("HASS_TOKEN") == "" {
		return
	}
	hass.SaveCons()
}

func handleNotifyConsumption() {
	hour, min, _ := time.Now().Clock()
	if !config.Conf.Enable.LogConsumption || !(hour == 10 && min == 0) {
		return
	}
	hass.GetYesterdayCons()
}

func handleMarketplace() {
	if !config.Conf.Enable.Marketplace {
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
	if !config.Conf.Enable.Bumphva {
		return
	}
	log.Println("Bump HVA")
	for _, user := range config.Conf.HaBump {
		for _, item := range user.Items {
			bumpha.Update(user.Identifier, item)
		}
	}
}

func stockWatcher() {
	if !config.Conf.Enable.Stockwatcher {
		return
	}
	for _, item := range config.Conf.StockWatcher {
		scraper.StockWatcher(&item)
	}
}

func queryArukereso() {
	if !config.Conf.Enable.Arukereso {
		return
	}
	log.Println("Arukereso")
	for _, item := range config.Conf.Arukereso {
		scraper.QueryArukereso(item.Url)
	}
}

func handleUpdateBB() {
	if !config.Conf.Enable.Bestbuy {
		return
	}
	scraper.UpdateBestBuy()
}

func handleSG() {
	if !config.Conf.Enable.Steamgifts || os.Getenv("SG_SESSID") == "" {
		return
	}
	log.Println("Sg")
	sg.QueryEntry()
}

func handleUpdateCovid() {
	if !config.Conf.Enable.Covid {
		return
	}
	log.Println("COVID")
	scraper.UpdateCovid()
}

func handleUpdateDHT() {
	if !config.Conf.Enable.Dht {
		return
	}
	log.Println("DHT")
	dht.ReadDHT(config.Conf.Dht.Pin)
}

func handleUpdateNcore() {
	if !config.Conf.Enable.Ncore || os.Getenv("NCORE_USERNAME") == "" || os.Getenv("NCORE_PASSWORD") == "" {
		return
	}
	log.Println("Ncore")
	scraper.UpdateNcore()
}

func handleUpdateFuel() {
	if !config.Conf.Enable.Fuel {
		return
	}
	log.Println("Fuel")
	sg.QueryEntry()
	scraper.UpdateFuelPrice()
}

func handleUpdateFixer() {
	if !config.Conf.Enable.Fixerio || os.Getenv("FIXERAPI") == "" {
		return
	}
	log.Println("Fixer")
	scraper.UpdateCurrencies()
}

func handleAwsCost() {
	if !config.Conf.Enable.Awscost {
		return
	}
	log.Println("AWS Cost")
	awscost.Update()
}

func handleBTC() {
	if !config.Conf.Enable.Btc || os.Getenv("BTC_API") == "" {
		return
	}
	log.Println("BTC")
	scraper.UpdateBTC()
}
