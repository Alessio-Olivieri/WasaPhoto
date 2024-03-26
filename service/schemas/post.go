package schemas

type Post struct {
	PostId     uint64    `json:"post_id"` // Unique identifier of the post (uint64 for compatibility with object_id)
	Text       string    `json:"content"` // Content of the post
	Username   string    `json:"username"`
	UserId     uint64    `json:"user_id"` // Unique identifier of the user (uint64 for compatibility with object_id)
	LikesCount int       `json:"likes_count"`
	Likes      []string  `json:"likes"`    // List of users who liked the post
	IsLiked    bool      `json:"is_liked"` // True if the requesting user liked the post
	Comments   []Comment `json:"comments"` // List of comments on the post
	Date       int64     `json:"date"`     // Date of the post
	Picture    []byte    `json:"image"`    // Reference to the photo schema (uint64 for compatibility with photo)
}

type Comment struct {
	CommentId uint64 `json:"comment_id"` // Unique identifier of the comment (uint64 for compatibility with object_id)
	PostId    uint64 `json:"post_id"`    // Unique identifier of the photo (uint64 for compatibility with object_id)
	Username  string `json:"username"`
	UserId    uint64 `json:"user_id"` // Unique identifier of the user who posted the comment (uint64 for compatibility with object_id)
	Text      string `json:"content"` // Content of the comment
	Date      int64  `json:"date"`    // Date of the comment
}
