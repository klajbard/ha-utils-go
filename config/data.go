package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Marketplace  MarketplaceConfig `yaml:"marketplace"`
	HaBump       []HaBumpConfig    `yaml:"habump"`
	Channels     []SlackChannel    `yaml:"channels"`
	StockWatcher []ItemStock       `yaml:"stockwatcher"`
	Arukereso    []Url             `yaml:"arukereso"`
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
	Bestbuy      bool `yaml:"bestbuy"`
	Stockwatcher bool `yaml:"stockwatcher"`
	Marketplace  bool `yaml:"marketplace"`
	Steamgifts   bool `yaml:"steamgifts"`
	Dht          bool `yaml:"dht"`
	Arukereso    bool `yaml:"arukereso"`
	Covid        bool `yaml:"covid"`
	Bumphva      bool `yaml:"bumphva"`
	Ncore        bool `yaml:"ncore"`
	Fuel         bool `yaml:"fuel"`
	Fixerio      bool `yaml:"fixerio"`
	Awscost      bool `yaml:"awscost"`
	Btc          bool `yaml:"btc"`
}

func (c *Configuration) GetConf() *Configuration {
	cfg := flag.String("cfg", "config.yaml", "config file path")
	flag.Parse()

	conf, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Println(err)
	}
	if err := yaml.Unmarshal(conf, c); err != nil {
		log.Println(err)
	}

	return c
}

var Conf Configuration

var Channels = map[string]string{}

func init() {
	Conf.GetConf()
	for _, channel := range Conf.Channels {
		Channels[channel.Name] = channel.Id
	}
}
