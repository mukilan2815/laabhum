package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Oms struct {
		BaseURL string `yaml:"baseURL"`
	} `yaml:"oms"`
	LogLevel     string `yaml:"log_level"`
	OMSAddress   string `yaml:"oms_address"`
	ServerAddress string `yaml:"server_address"`
}

func LoadConfig() *Config {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return &cfg
}
