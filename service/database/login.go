package database

import (
	"errors"
	"fmt"
)

// Exists_users returns the id of the user if the user exists in the database.
// Returns 18446744073709551615 if the user does not exist
func (db *appdbimpl) Login_db(username string) (uint64, error) {
	rows, err := db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	defer func() { _ = rows.Close() }()
	if err != nil {
		return 18446744073709551615, fmt.Errorf("error querying database: %w", err)
	}

	var user_id uint64
	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return 18446744073709551615, fmt.Errorf("error reading rows: %w", err)
		}
		// User exists
		return user_id, nil
	}

	_, err = db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
	if err != nil {
		return 18446744073709551615, fmt.Errorf("error inserting new user: %w", err)
	}

	rows, err = db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	if err != nil {
		return 18446744073709551615, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return 18446744073709551615, fmt.Errorf("error reading rows: %w", err)
		}
		// User exists
		return user_id, nil
	}

	return 18446744073709551615, errors.New("error inserting new user")
}
