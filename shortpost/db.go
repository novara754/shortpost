package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// MustOpenSQL creates a connection using the given Postgres URL.
// Panics on error.
func MustOpenSQL(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Panicf("Failed to open database connection: %s", err.Error())
	}
	return db
}
