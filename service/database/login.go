package database

import (
	"database/sql"
	"errors"
)

// Exists_users returns the id of the user if the user exists in the database.
func (db *appdbimpl) Login_db(username string) (uint64, error) {
	var user_id uint64
	// search user in the database
	err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?;", username).Scan(&user_id)
	if errors.Is(err, sql.ErrNoRows) {
		// if the user does not exist, create a new user
		_, err = db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
		if err != nil {
			return 18446744073709551615, err
		}
		// return the id of the new user
		err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?;", username).Scan(&user_id)
		if errors.Is(err, sql.ErrNoRows) {
			return 18446744073709551615, err
		}
		if err != nil {
			return 18446744073709551615, err
		}
		return user_id, nil
	} else if err == nil {
		// if the user exists, return the id of the user
		return user_id, nil
	}
	// return an error if the query fails
	return 18446744073709551615, err
}
