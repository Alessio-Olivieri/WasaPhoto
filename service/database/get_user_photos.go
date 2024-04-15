package database

import (
	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) Get_user_photos(user_id_request uint64, user_id_profile uint64, page_number int) ([]schemas.Post, int, error) {
	var postList []schemas.Post
	var totalPosts int
	// retrieve the total number of posts
	err := db.c.QueryRow(`SELECT count(*) FROM 
		photos INNER JOIN users ON users.user_id = photos.user_id
		WHERE 
		photos.user_id IN (SELECT followed_id FROM Followers WHERE follower_id = ?)`, user_id_request).Scan(&totalPosts)
	if err != nil {
		return postList, 0, err
	}

	TotalPages := totalPosts / 10
	if totalPosts < page_number*10 {
		return postList, 0, ErrPageNumberOutOfBound
	}

	rows, err := db.c.Query(`SELECT Photos.photo_id, 
		Photos.image,
		Photos.user_id, 
		(SELECT username FROM Users WHERE user_id = ?),
		Photos.text,
		Photos.date,
		(SELECT EXISTS(SELECT TRUE FROM Likes WHERE Photos.photo_id = photo_id AND user_id = ?)) AS isLiked 
		FROM Photos WHERE user_id = ?
		ORDER BY Photos.date DESC
		LIMIT 10 OFFSET ?`, user_id_request, user_id_request, user_id_profile, page_number)
	if err != nil {
		return postList, 0, err
	}
	defer rows.Close()

	postList, err = db.retrievePosts(rows)
	if err != nil {
		return postList, 0, err
	}

	return postList, TotalPages, nil
}
