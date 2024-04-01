package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) DeletePost(post_id uint64, user_id uint64) error {
	//check if user is authorized
	var authorized bool
	var exists bool

	//check if photo exists
	err := db.c.QueryRow("SELECT true FROM Photos WHERE photo_id = ?", post_id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPhotoNotExists
	}
	if err != nil {
		return err
	}

	err = db.c.QueryRow("SELECT true FROM Photos WHERE user_id = ? and photo_id = ?", user_id, post_id).Scan(&authorized)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrForbidden
	}
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`DELETE from Photos WHERE photo_id = ?`, post_id)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(`DELETE from Comments WHERE photo_id = ?`, post_id)
	if err != nil {
		return err
	}
	return nil
}
