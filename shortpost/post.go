package main

import (
	"database/sql"
)

// Post represents a post made and displayed on the
// shortpost website.
type Post struct {
	ID      int64
	Author  User
	Content string
}

// CreatePostTable creates the database table for posts
// if it doesnt already exist.
func CreatePostTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS posts (
	id serial PRIMARY KEY,
	content varchar(240) NOT NULL,
	author_id serial NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (author_id) REFERENCES users
);
	`)
	return err
}

// GetAllPost queries the database to retrieve all
// posts.
func GetAllPost(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
SELECT posts.id, content, author_id, username
FROM posts
INNER JOIN users ON users.id = author_id
ORDER BY posts.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.ID, &post.Content, &post.Author.ID, &post.Author.Username); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// InsertPost inserts a new post into the database with the given
// author and content.
func InsertPost(authorID int64, content string) error {
	_, err := db.Exec(`
INSERT INTO posts (author_id, content) VALUES ($1, $2);
	`, authorID, content)
	return err
}
