package database

import (
	"database/sql"
	"fmt"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (db *appdbimpl) retrievePosts(rows *sql.Rows) ([]schemas.Post, error) {
	var postList []schemas.Post

	// create each post
	for rows.Next() {
		var post schemas.Post
		err := rows.Scan(&post.PostId, &post.Picture, &post.UserId, &post.Username, &post.Text, &post.Date, &post.IsLiked)
		if err != nil {
			return nil, err
		}
		fmt.Println("DATABASE: post retrieved", post.PostId, post.Username, post.UserId, post.Text, post.Date, post.IsLiked)

		// get the usernames of who putted like
		post.Likes, err = db.GetLikes(post.PostId)
		if err != nil {
			return nil, err
		}
		post.LikesCount = len(post.Likes)

		// get the comments of a picture
		post.Comments, err = db.GetComments(post.PostId)
		if err != nil {
			return nil, err
		}

		postList = append(postList, post)
	}
	return postList, nil
}
