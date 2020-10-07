package main

import (
	"log"
)

var config = MustLoadConfig()
var db = MustOpenSQL(config.PostgresURL)

func main() {
	if err := CreatePostTable(db); err != nil {
		log.Fatalf("Failed to create posts table: %s", err.Error())
	}
}
