package database

import (
	"errors"
	"log"
)

// Exists_users returns the id of the user if the user exists in the database.
// Returns -1 if the user does not exist
func (db *appdbimpl) Exists_user(username string) (int, error) {
	log.Println("checking if user exists")
	rows, err := db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	defer func() { _ = rows.Close() }()
	if err != nil {
		return -1, err
	}
	var user_id int
	if rows.Next() {
		if err := rows.Scan(&user_id); err != nil {
			return -1, err
		}
		return user_id, nil
	}

	return -1, errors.New("error checking if user exists")
}
