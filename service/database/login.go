package database

// Exists_users returns the id of the user if the user exists in the database.
// Returns 18446744073709551615 if the user does not exist
func (db *appdbimpl) Login_db(username string) (uint64, error) {
	rows, err := db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	if err != nil {
		return 18446744073709551615, err
	}
	defer rows.Close()

	var user_id uint64
	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return 18446744073709551615, err
		}
		// User exists
		return user_id, nil
	}

	_, err = db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
	if err != nil {
		return 18446744073709551615, err
	}

	rows, err = db.c.Query("SELECT user_id FROM Users WHERE username = ?;", username)
	if err != nil {
		return 18446744073709551615, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return 18446744073709551615, err
		}
		// User exists
		return user_id, nil
	}

	return 18446744073709551615, err
}
