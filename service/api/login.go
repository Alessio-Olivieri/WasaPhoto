package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("login called!")

	// Get the username and password from the request
	username := "Aleksandra"

	// Check if the username exists in the database
	// The database returns the unique id of the user
	uid, err := rt.db.Exists_user(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("login error:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//the database returns a unique session id corresponding to the userID
	sessionID, err := rt.db.Add_session(uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("error creating session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("session created")
	log.Println(sessionID)

	// Create a cookie that will be used to authenticate the user
	log.Println("Creating cookie")
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: fmt.Sprintf("%d", sessionID),
		Path:  "/",
	})
	fmt.Fprintf(w, "cookie set: %d", sessionID)

	http.Redirect(w, r, "/context", http.StatusSeeOther)

}
