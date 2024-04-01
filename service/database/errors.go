package database

import "errors"

var ErrForbidden = errors.New("unauthorized user")
var ErrBanned = errors.New("user is banned")

var ErrCommentNotExists = errors.New("comment not in the database")

var ErrPhotoNotExists = errors.New("photo not in the database")

var ErrUserNotExists = errors.New("user not in the database")
var ErrUsernameTaken = errors.New("username already exists")

var ErrEmptyStream = errors.New("stream is empty")
var ErrPageNumberOutOfBound = errors.New("page number asked is too high")

var ErrBanNotExists = errors.New("ban not in the database")
var ErrBanAlreadyExists = errors.New("ban already exists")

var ErrAlreadyLiked = errors.New("photo already liked")
var ErrLikeNotExists = errors.New("photo already unliked")

var ErrAlreadyFollowed = errors.New("user already followed")
var ErrFollowNotExists = errors.New("user already unfollowed")
