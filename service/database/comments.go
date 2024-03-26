package database

import (
	"time"

	"github.com/Alessio-Olivieri/wasaProject/service/schemas"
)

func (db *appdbimpl) PostComment(photo_id uint64, user_id uint64, text string) error {
	time := time.Now().Unix()

	_, err := db.c.Exec(`INSERT into Comments (photo_id, user_id, text, date) VALUES (?, ?, ?, ?)`, photo_id, user_id, text, time)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteComment(commentId uint64) error {
	_, err := db.c.Exec(`DELETE from Comments WHERE comment_id = ?`, commentId)
	if err != nil {
		return err
	}
	return nil
}

// get the comments of a picture
func (db *appdbimpl) GetComments(photo_id uint64) ([]schemas.Comment, error) {
	var comments []schemas.Comment
	rows, err := db.c.Query(`SELECT comment_id, user_id, text, date FROM Comments WHERE photo_id = ?`, photo_id)
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
