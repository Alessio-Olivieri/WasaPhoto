package database

import (
	"fmt"
	"log"
)

func (db *appdbimpl) Add_session(user_id int) (int, error) {
	log.Println("Creating session for ", user_id)

	// Create a new session for the user specified by user_id
	_, err := db.c.Exec(`INSERT INTO Sessions ('user_id') VALUES (?);`, user_id)
	if err != nil {
		return -1, fmt.Errorf("error creating session: %w", err)
	}

	// Get the session id
	var session_id int
	err = db.c.QueryRow(`SELECT session_id FROM Sessions WHERE user_id = ?;`, user_id).Scan(&session_id)
	if err != nil {
		return -1, fmt.Errorf("error getting session id: %w", err)
	}
	log.Println("session created:", session_id)
	return session_id, nil
}

func (db *appdbimpl) Get_userID(session_id int) (int, error) {
	var user_id int
	err := db.c.QueryRow(`SELECT user_id FROM Sessions WHERE session_id = ?;`, session_id).Scan(&user_id)
	if err != nil {
		return -1, fmt.Errorf("error getting user id from session: %w", err)
	}
	return user_id, nil
}
