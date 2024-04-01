package database

import (
	"database/sql"
	"errors"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) Get_stream(request_user_id uint64, page int) (schemas.Stream, error) {
	var stream schemas.Stream
	var totalPosts int
	//retrieve the total number of posts
	err := db.c.QueryRow(`SELECT count(*) FROM 
		photos INNER JOIN users ON users.user_id = photos.user_id
		WHERE 
		photos.user_id IN (SELECT followed_id FROM Followers WHERE follower_id = ?)`, request_user_id).Scan(&totalPosts)
	if err != nil {
		return stream, err
	}
	if totalPosts == 0 {
		return stream, ErrEmptyStream
	}

	stream.TotalPages = int(totalPosts / 10)
	if totalPosts < page*10 {
		return stream, ErrPageNumberOutOfBound
	}

	rows, err := db.c.Query(`
		 SELECT 
			Photos.photo_id, 
			Photos.image,
			Photos.user_id, Users.username,
			Photos.text,
			Photos.date,
			(SELECT EXISTS(SELECT TRUE FROM Likes WHERE Photos.photo_id = photo_id AND user_id = ?)) AS isLiked
		FROM 
			photos INNER JOIN users ON users.user_id = photos.user_id
		WHERE 
			photos.user_id IN (SELECT followed_id FROM Followers WHERE follower_id = ?)
		ORDER BY Photos.date DESC
		LIMIT 10 OFFSET ?
		`, request_user_id, request_user_id, page*10)
	if errors.Is(err, sql.ErrNoRows) {
		return stream, ErrPageNumberOutOfBound
	}
	if err != nil {
		return stream, err
	}

	stream.Posts, err = db.retrievePosts(rows)
	if err != nil {
		return stream, err
	}

	return stream, nil
}
