package main

import (
	"database/sql"
)

// Post represents a post made and displayed on the
// shortpost website.
type Post struct {
	ID         int64
	AuthorName string
	Content    string
}

// CreatePostTable creates the database table for posts
// if it doesnt already exist.
func CreatePostTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS posts (
	id serial PRIMARY KEY,
	content varchar(240) NOT NULL,
	author_name varchar(30) NOT NULL
);
	`)
	return err
}
