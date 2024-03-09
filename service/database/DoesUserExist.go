package database

import (
	"errors"
	"fmt"
)

// Exists_users returns the id of the user if the user exists in the database.
// Returns -1 if the user does not exist
func (db *appdbimpl) Exists_user(username string) (int, error) {
	rows, err := db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	defer func() { _ = rows.Close() }()
	if err != nil {
		return -1, fmt.Errorf("error querying database: %w", err)
	}

	var user_id int
	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return -1, fmt.Errorf("error reading rows: %w", err)
		}
		// User exists
		return user_id, nil
	}

	_, err = db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
	if err != nil {
		return -1, fmt.Errorf("error inserting new user: %w", err)
	}

	rows, err = db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return -1, fmt.Errorf("error reading rows: %w", err)
		}
		// User exists
		return user_id, nil
	}

	return -1, errors.New("error inserting new user")
}
