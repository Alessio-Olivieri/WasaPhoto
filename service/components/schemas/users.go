package schemas

type UserList struct {
	Users []string `json:"users"`
}

type Username struct {
	Username string `json:"username"`
}

type GetUserProfile struct {
	Followers       []string `json:"followers"`
	Followers_count int      `json:"followers_count"`
	IsBanned        bool     `json:"is_banned"`
	IsFollowing     bool     `json:"is_following"`
	Posts           []Post   `json:"posts"`
}
