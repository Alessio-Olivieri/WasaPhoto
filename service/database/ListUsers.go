package database

import (
	"log"
)

// return a list with the name of all the users
func (db *appdbimpl) ListUsers() ([]string, error) {
	log.Println(" listing users")
	rows, err := db.c.Query(`SELECT username FROM Users;`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var ret []string
	var username string
	log.Println("enering loop")

	for rows.Next() {
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		ret = append(ret, username)
	}
	log.Println("returning list of users, length:", len(ret))
	return ret, nil
}