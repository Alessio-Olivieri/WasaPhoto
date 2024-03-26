package schemas

type UserList struct {
	Users []string `json:"users"`
}

type Username struct {
	Username string `json:"username"`
}

type GetUserProfile struct {
	Username        string   `json:"username"`
	UserId          uint64   `json:"user_id"`
	Followers       []string `json:"followers"`
	Followers_count int      `json:"followers_count"`
	Posts           []Post   `json:"posts"`
	IsBanned        bool     `json:"isBanned"`
	IsFollowing     bool     `json:"isFollowing"`
}
