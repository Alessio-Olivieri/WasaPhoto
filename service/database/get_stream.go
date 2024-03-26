package database

import (
	"github.com/Alessio-Olivieri/wasaProject/service/schemas"
)

func (db *appdbimpl) Get_stream(request_user_id uint64, page int) ([]schemas.Post, error) {
	var postList []schemas.Post

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
		`, request_user_id, request_user_id, page)
	if err != nil {
		return postList, err
	}

	//create each post
	for rows.Next() {
		var post schemas.Post
		err = rows.Scan(&post.PostId, &post.Picture, &post.UserId, &post.Username, &post.Text, &post.Date, &post.IsLiked)
		if err != nil {
			return nil, err
		}
		//determine if the picture il liked by the requesting user
		post.IsLiked, err = db.IsLiked(post.PostId, request_user_id)
		if err != nil {
			return nil, err
		}

		//get the usernames of who putted like
		post.Likes, err = db.GetLikes(post.PostId)
		if err != nil {
			return nil, err
		}
		post.LikesCount = len(post.Likes)

		//get the comments of a picture
		post.Comments, err = db.GetComments(post.PostId)
		if err != nil {
			return nil, err
		}

		postList = append(postList, post)
	}

	return postList, nil
}
