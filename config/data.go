package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/klajbard/ha-utils-go/types"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Marketplace  MarketplaceConfig `yaml:"marketplace"`
	HaBump       []HaBumpConfig    `yaml:"habump"`
	Channels     []SlackChannel    `yaml:"channels"`
	StockWatcher []ItemStock       `yaml:"stockwatcher"`
	Arukereso    []Url             `yaml:"arukereso"`
	Dht          DhtConfig         `yaml:"dht"`
	Enable       EnableConfig      `yaml:"enable"`
	Silence      bool              `yaml:"silence"`
}

type Url struct {
	Url string `yaml:"url"`
}

type MarketplaceConfig struct {
	Jofogas     []MarketplaceName `yaml:"jofogas"`
	Hardverapro []MarketplaceName `yaml:"hardverapro"`
}

type HaBumpConfig struct {
	Identifier string   `yaml:"identifier"`
	Items      []HaItem `yaml:"items"`
}

type DhtConfig struct {
	Pin int `yaml:"pin"`
}

type HaItem struct {
	Name  string `yaml:"name"`
	Id    string `yaml:"id"`
	Hour  int    `yaml:"hour"`
	Start int    `yaml:"start"`
}

type MarketplaceName struct {
	Name string `yaml:"name"`
}

type SlackChannel struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type ItemStock struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Query string `yaml:"query"`
}

type EnableConfig struct {
	Bestbuy        bool `yaml:"bestbuy"`
	Stockwatcher   bool `yaml:"stockwatcher"`
	Marketplace    bool `yaml:"marketplace"`
	Steamgifts     bool `yaml:"steamgifts"`
	Dht            bool `yaml:"dht"`
	Arukereso      bool `yaml:"arukereso"`
	Covid          bool `yaml:"covid"`
	Bumphva        bool `yaml:"bumphva"`
	Ncore          bool `yaml:"ncore"`
	Fuel           bool `yaml:"fuel"`
	Fixerio        bool `yaml:"fixerio"`
	Awscost        bool `yaml:"awscost"`
	Btc            bool `yaml:"btc"`
	LogConsumption bool `yaml:"logcons"`
}

var cfg = flag.String("cfg", "config.yaml", "config file path")
var Conf Configuration
var Channels = map[string]string{}

func (c *Configuration) GetConf() *Configuration {
	var conf []byte
	resp, err := http.Get(fmt.Sprintf("https://%s@cdn.klajbar.com/conf/config.yaml", os.Getenv("CDN_CRED")))
	if err != nil {
		log.Println(err)
		flag.Parse()
		conf, err = ioutil.ReadFile(*cfg)
		if err != nil {
			log.Println(err)
		}

	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		configData := types.AWSResponse{}

		err = json.Unmarshal([]byte(string(respBody)), &configData)
		if err != nil {
			log.Fatal(err)
		}
		conf = configData.Body.Data
	}

	if err := yaml.Unmarshal(conf, c); err != nil {
		log.Println(err)
	}

	for _, channel := range c.Channels {
		Channels[channel.Name] = channel.Id
	}

	return c
}
