package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) PostComment(photo_id uint64, user_id uint64, text string) (schemas.Comment, error) {

	time := time.Now().Unix()

	var comment schemas.Comment

	var user_id_photo uint64
	err := db.c.QueryRow(`SELECT user_id FROM Photos WHERE photo_id = ?`, photo_id).Scan(&user_id_photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comment, ErrPhotoNotExists
		}
		return comment, err
	}

	banned, err := db.IsBanned(user_id_photo, user_id)
	if err != nil {
		return comment, err
	}
	if banned {
		return comment, ErrBanned
	}

	result, err := db.c.Exec(`INSERT into Comments (photo_id, user_id, text, date) VALUES (?, ?, ?, ?)`, photo_id, user_id, text, time)
	if err != nil {
		return comment, err
	}

	username, err := db.Get_username_from_userId(user_id)
	if err != nil {
		return comment, err
	}
	comment_id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment = schemas.Comment{UserId: user_id, Text: text, Date: time, PostId: photo_id, Username: username, CommentId: uint64(comment_id)}

	return comment, nil
}

func (db *appdbimpl) DeleteComment(comment_id uint64, user_id uint64) error {
	// check if user is authorized
	var authorized bool
	var exists bool

	// check if comment exists
	err := db.c.QueryRow("SELECT true FROM Comments WHERE comment_id = ?", comment_id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCommentNotExists
		}
		return err
	}

	err = db.c.QueryRow("SELECT true FROM Comments WHERE user_id = ? and comment_id = ?", user_id, comment_id).Scan(&authorized)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrForbidden
		}
		return err
	}

	_, err = db.c.Exec(`DELETE from Comments WHERE comment_id = ?`, comment_id)
	if err != nil {
		return err
	}
	return nil
}

// get the comments of a picture
func (db *appdbimpl) GetComments(photo_id uint64) ([]schemas.Comment, error) {
	var comments []schemas.Comment
	rows, err := db.c.Query(`SELECT comment_id, user_id, text, date FROM Comments WHERE photo_id = ?`, photo_id)
	if errors.Is(err, sql.ErrNoRows) {
		return comments, ErrPhotoNotExists
	}
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment schemas.Comment
		err = rows.Scan(&comment.CommentId, &comment.UserId, &comment.Text, &comment.Date)
		if err != nil {
			return nil, err
		}

		comment.Username, err = db.Get_username_from_userId(comment.UserId)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
