package database

import (
	"strconv"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) Get_userId_from_username(username string) (uint64, error) {
	var user_id_s string
	err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?", username).Scan(&user_id_s)
	user_id, err := strconv.ParseUint(user_id_s, 10, 64)
	if err != nil {
		return 18446744073709551615, err
	}
	return user_id, err
}
