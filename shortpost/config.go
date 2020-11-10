package main

import (
	"log"
	"os"
)

// Config is a structure to contain configuration
// for the program.
type Config struct {
	Port        string
	PostgresURL string
}

// MustLoadConfig loads configuration data from environment variables.
// Panics on missing values.
func MustLoadConfig() Config {
	config := Config{
		Port:        os.Getenv("PORT"),
		PostgresURL: os.Getenv("POSTGRES_URL"),
	}

	if config.Port == "" {
		log.Panicf("Empty PORT variable")
	}
	if config.PostgresURL == "" {
		log.Panicf("Empty POSTGRES_URL variable")
	}

	return config
}
