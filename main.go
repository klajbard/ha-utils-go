package main

import (
	"fmt"
	"net/http"

	"./bumpha"
	"./consumptions"
	"./scraper"
	"./sg"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()
	mux.GET("/cons/:device", consumptions.GetAllCons)
	mux.GET("/cons/:device/:date", consumptions.GetCons)
	mux.PUT("/cons/:device", consumptions.PutCons)
	mux.GET("/scrape/pico", scraper.PicoScraper)
	mux.GET("/scrape/fuel", scraper.GetFuel)
	mux.GET("/jofogas/:item", scraper.GetJofogas)
	mux.GET("/hvapro/:item", scraper.GetHvapro)
	mux.GET("/sg", sg.QueryEntry)
	mux.GET("/ncore", scraper.Ncore)
	mux.GET("/bumpha/:fid/:name", bumpha.BumpHa)

	mux.POST("/test", testHandler)

	http.ListenAndServe(":5500", mux)
}

func testHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	fmt.Println(r.Form)
}
