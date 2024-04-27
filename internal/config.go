package internal

import (
	"github.com/go-yaml/yaml"
	"log"
	"os"
)

type Conversation struct {
	Type       string `yaml:"type"`
	FacebookID string `yaml:"facebook_id"`
	SignalID   string `yaml:"signal_id"`
}

type Person struct {
	SignalNumber string `yaml:"signal_number"`
	FacebookID   string `yaml:"facebook_id"`
}

type Config struct {
	Entrypoint    string         `yaml:"entrypoint"`
	People        []Person       `yaml:"people"`
	Conversations []Conversation `yaml:"conversations"`
}

var LoadedConfig *Config

func GetConfig() *Config {
	if LoadedConfig == nil {
		LoadConfig()
	}

	return LoadedConfig
}

func LoadConfig() {
	log.Printf("Loading config")

	var config Config

	data, err := os.ReadFile("config.yml")

	if err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.Fatalf("Can't parse config: %v", err)
	}

	LoadedConfig = &config
}
