package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetSearchUser(prompt string, requestingUser uint64) ([]string, error) {
	/* returns at most 24 ReducedUsers records with the requested prompt
	in their username among those who have not banned the requesteduser */
	// Implements GET /users/

	var result []string
	rows, err := db.c.Query(`
		SELECT username FROM users
		WHERE username LIKE '%' || ? || '%' AND
		NOT EXISTS (SELECT 1 FROM bans WHERE
			banner_id = Users.user_id AND banned_id = ?)
		ORDER BY LENGTH(username) ASC
		LIMIT 20;`, prompt, requestingUser)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotExists
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		result = append(result, username)
	}
	return result, nil
}
