package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//login routes
	rt.router.GET("/login", rt.login)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	//user
	rt.router.GET("/users/", rt.get_users_list)

	return rt.router
}
