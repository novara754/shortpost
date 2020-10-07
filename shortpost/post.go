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
	author_name varchar(30) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
	`)
	return err
}

// GetAllPost queries the database to retrieve all
// posts.
func GetAllPost(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
SELECT id, content, author_name FROM posts
ORDER BY created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.ID, &post.Content, &post.AuthorName); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// InsertPost inserts a new post into the database with the given
// author and content.
func InsertPost(author string, content string) error {
	_, err := db.Exec(`
INSERT INTO posts (content, author_name) VALUES ($1, $2);
	`, content, author)
	return err
}
