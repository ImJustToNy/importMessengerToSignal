package internal

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conversation struct {
	Type          string `yaml:"type"`
	FacebookID    string `yaml:"facebook_id"`
	SignalID      string `yaml:"signal_id"`
	FromTimestamp int64  `yaml:"resume_from_timestamp"`
}

type Person struct {
	SignalNumber string `yaml:"signal_number"`
	FacebookID   string `yaml:"facebook_id"`
}

type Config struct {
	Entrypoint       string         `yaml:"entrypoint"`
	People           []Person       `yaml:"people"`
	Conversations    []Conversation `yaml:"conversations"`
	RestartDbusEvery int            `yaml:"restart_dbus_every"`
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
