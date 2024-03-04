package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("login called!")

	// Get the username and password from the request

	username := "Aleksandra"

	rt.db.Exists_user(username)

}
