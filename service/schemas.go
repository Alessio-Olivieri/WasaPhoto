package service

type PhotoListElement struct {
	PhotoId     uint64 `json:"photoId"`
	ReleaseDate string `json:"date"`
	Caption     string `json:"caption"`
	PublisherId uint64 `json:"userId"`
	Likes       int    `json:"likecount"`
}
