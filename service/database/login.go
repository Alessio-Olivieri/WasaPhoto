package database

import (
	"database/sql"
	"errors"
)

// Exists_users returns the id of the user if the user exists in the database.
func (db *appdbimpl) Login_db(username string) (uint64, bool, error) {
	var user_id uint64
	// search user in the database
	err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?;", username).Scan(&user_id)
	if errors.Is(err, sql.ErrNoRows) {
		// if the user does not exist, create a new user
		_, err = db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
		if err != nil {
			return 18446744073709551615, false, err
		}
		// return the id of the new user
		err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?;", username).Scan(&user_id)
		if err != nil {
			return 18446744073709551615, false, err
		}
		return user_id, true, nil
	} else if err == nil {
		// if the user exists, return the id of the user
		return user_id, false, nil
	}
	// return an error if the query fails
	return 18446744073709551615, false, err
}
