package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config is a structure to contain configuration
// for the program.
type Config struct {
	PostgresURL string `json:"postgresURL"`
}

// MustLoadConfig loads configuration data from `config.json`.
// Panics on error.
func MustLoadConfig() (config *Config) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panicf("Failed to read config: %s", err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Panicf("Failed to parse config: %s", err.Error())
	}

	return
}
