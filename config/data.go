package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Marketplace MarketplaceConfig `yaml:"marketplace"`
	HaBump      []HaBumpConfig    `yaml:"habump"`
	Channels    []SlackChannel    `yaml:"channels"`
}

type MarketplaceConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Jofogas     []MarketplaceName `yaml:"jofogas"`
	Hardverapro []MarketplaceName `yaml:"hardverapro"`
}

type HaBumpConfig struct {
	Identifier string   `yaml:"identifier"`
	Items      []HaItem `yaml:"items"`
}

type HaItem struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type MarketplaceName struct {
	Name string `yaml:"name"`
}

type SlackChannel struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
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
