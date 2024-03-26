package database

import (
	"github.com/Alessio-Olivieri/wasaProject/service/schemas"
)

func (db *appdbimpl) Get_user_photos(user_id_request uint64, user_id_profile uint64) ([]schemas.Post, error) {
	var postlist []schemas.Post
	rows, err := db.c.Query(`SELECT photo_id, text, like_count, image, "date" FROM Photos WHERE user_id = ?;`, user_id_profile)
	if err != nil {
		return postlist, err
	}

	//create each post
	for rows.Next() {
		var post schemas.Post
		err = rows.Scan(&post.PostId, &post.Text, &post.LikesCount, &post.Picture, &post.Date)
		if err != nil {
			return nil, err
		}
		//determine if the picture il liked by the requesting user
		post.IsLiked, err = db.IsLiked(post.PostId, user_id_request)
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

		postlist = append(postlist, post)
	}

	return postlist, nil
}
