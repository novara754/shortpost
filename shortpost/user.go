package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user who can make new posts.
type User struct {
	ID       int64
	Username string
}

// CreateUserTable creates the database table for posts
// if it doesnt already exist.
func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	username varchar(30) NOT NULL,
	password varchar(60) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
	`)
	return err
}

// GetUserByID queries the database for a user based
// on the given user id.
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	rows := db.QueryRow(`
SELECT id, username FROM users WHERE id = $1;
	`, id)

	user := User{}
	if err := rows.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}

	return &user, nil
}

// InsertUser inserts a new user into the database with the given
// username and password.
// The password is hashed before being stored.
func InsertUser(username string, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, username;
	`, username, passwordHash)

	user := User{}
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser verifies the given username and password combo and returns the
// matching user.
func AuthenticateUser(username string, password string) (*User, error) {
	row := db.QueryRow(`
SELECT id, username, password
FROM users
WHERE username = $1
	`, username)

	user := User{}
	var passwordHash []byte
	if err := row.Scan(&user.ID, &user.Username, &passwordHash); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}
