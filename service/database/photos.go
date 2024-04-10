package database

import (
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) Make_photo(user_id uint64, caption string, picture multipart.File) (schemas.Post, error) {
	// Get current time
	time := time.Now().Unix()

	var fileBytes []byte
	var err error

	if picture != nil {
		fileBytes, err = io.ReadAll(picture)
		if err != nil {
			return schemas.Post{}, err
		}
	}

	_, err = db.c.Exec("INSERT INTO Photos (user_id, text, image, date) VALUES (?, ?, ?, ?)", user_id, caption, fileBytes, time)
	if err != nil {
		return schemas.Post{}, err
	}

	// get username for the response
	username, err := db.Get_username_from_userId(user_id)
	if err != nil {
		return schemas.Post{}, err
	}

	response := schemas.Post{UserId: user_id, Text: caption, Picture: fileBytes, Date: time, Username: username}
	return response, err
}

func (db *appdbimpl) Exists_photo(photo_id uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT 1 FROM Photos WHERE photo_id = ?", photo_id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return exists, err
	}
	return true, nil
}
