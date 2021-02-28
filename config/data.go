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
}

type MarketplaceConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Jofogas     []MarketplaceName `yaml:"jofogas"`
	Hardverapro []MarketplaceName `yaml:"hardverapro"`
}

type HaBumpConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type MarketplaceName struct {
	Name string `yaml:"name"`
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

func init() {
	Conf.GetConf()
}
