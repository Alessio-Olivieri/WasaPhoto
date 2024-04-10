package database

import (
	"database/sql"
	"errors"
	"strconv"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) Get_userId_from_username(username string) (uint64, error) {
	// check if user exists
	exists, err := db.Exists_user(username)
	if err != nil {
		return 18446744073709551615, err
	}
	if !exists {
		return 18446744073709551615, ErrUserNotExists
	}

	var user_id_s string
	err = db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?", username).Scan(&user_id_s)
	if errors.Is(err, sql.ErrNoRows) {
		return 18446744073709551615, ErrUserNotExists
	}
	if err != nil {
		return 18446744073709551615, err
	}
	user_id, err := strconv.ParseUint(user_id_s, 10, 64)
	if err != nil {
		return 18446744073709551615, err
	}
	return user_id, err
}

// GetName is an example that shows you how to query data
func (db *appdbimpl) Get_username_from_userId(userId uint64) (string, error) {
	// check if user exists
	exists, err := db.Exists_user(userId)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserNotExists
	}

	var username string
	err = db.c.QueryRow("SELECT username FROM Users WHERE user_id = ?", userId).Scan(&username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUserNotExists
	}
	if err != nil {
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) Exists_user(identifier interface{}) (bool, error) {
	var exists bool
	var err error

	switch v := identifier.(type) {
	case string:
		// Identifier is a username
		err = db.c.QueryRow("SELECT count(*) FROM Users WHERE username = ?", v).Scan(&exists)
		if err != nil {
			return exists, err
		}
	case uint64:
		// Identifier is a user ID
		err = db.c.QueryRow("SELECT count(*) FROM Users WHERE user_id = ?", v).Scan(&exists)
		if err != nil {
			return exists, err
		}
	default:
		return false, errors.New("Exists_user: invalid identifier type")
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return exists, err
}
