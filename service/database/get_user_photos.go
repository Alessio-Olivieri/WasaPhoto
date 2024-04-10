package database

import (
	"database/sql"
	"errors"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) Get_user_photos(user_id_request uint64, user_id_profile uint64) ([]schemas.Post, error) {
	var postList []schemas.Post
	rows, err := db.c.Query(`SELECT Photos.photo_id, 
		Photos.image,
		Photos.user_id, 
		(SELECT username FROM Users WHERE user_id = ?),
		Photos.text,
		Photos.date,
		(SELECT EXISTS(SELECT TRUE FROM Likes WHERE Photos.photo_id = photo_id AND user_id = ?)) AS isLiked 
		FROM Photos WHERE user_id = ?;`, user_id_request, user_id_request, user_id_profile)
	if errors.Is(err, sql.ErrNoRows) {
		return postList, ErrUserNotExists
	}
	if err != nil {
		return postList, err
	}
	defer rows.Close()

	postList, err = db.retrievePosts(rows)
	if err != nil {
		return postList, err
	}

	return postList, nil
}
