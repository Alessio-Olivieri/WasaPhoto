package api

import (
	"errors"
	"net/http"
)

var ErrUserNotExists = errors.New("username parameter missing in request path")

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// login routes
	rt.router.POST("/login", rt.wrap(rt.login, false))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// stream
	rt.router.GET("/stream", rt.wrap(rt.get_stream, true))

	// user
	rt.router.GET("/users/", rt.wrap(rt.search_users, true))
	rt.router.GET("/users/:username", rt.wrap(rt.GetUserProfile, true))
	rt.router.PUT("/settings/username", rt.wrap(rt.changeUsername, true))

	// follower
	rt.router.PUT("/followed/:username", rt.wrap(rt.put_follower, true))
	rt.router.DELETE("/followed/:username", rt.wrap(rt.delete_follower, true))

	// ban
	rt.router.PUT("/banned/:username", rt.wrap(rt.put_ban, true))
	rt.router.DELETE("/banned/:username", rt.wrap(rt.delete_ban, true))

	// like
	rt.router.PUT("/photos/:photo_id/likes/me", rt.wrap(rt.put_like, true))
	rt.router.DELETE("/photos/:photo_id/likes/me", rt.wrap(rt.delete_like, true))
	// rt.router.Get("/photos/:photo_id/likesAmount", rt.wrap(rt.get_likes, true))

	// comments
	rt.router.POST("/photos/:photo_id/comments/", rt.wrap(rt.post_comment, true))
	rt.router.DELETE("/photos/:photo_id/comments/:comment_id", rt.wrap(rt.delete_comment, true))

	// photo
	rt.router.POST("/photos/", rt.wrap(rt.post_photo, true))
	rt.router.DELETE("/photos/:photo_id", rt.wrap(rt.delete_photo, true))

	return rt.router
}
