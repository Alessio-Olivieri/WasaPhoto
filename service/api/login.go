package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// login is the handler for the POST /login route
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Login:\n"

	username := r.URL.Query().Get("username")

	// Check if the user exists in the database, if it does not exist, create it
	userID, created, err := rt.db.Login_db(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error(message + "error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Authenticated user " + username + " with ID " + strconv.FormatUint(userID, 10) + "\n"

	// response contains the user ID associated to the username
	response := map[string]uint64{"userId": userID}

	// Send the response to the client
	w.Header().Set("Content-Type", "application/json")
	if created {
		w.WriteHeader(http.StatusCreated)
		message = message + "User created\n"
	} else {
		w.WriteHeader(http.StatusOK)
		message = message + "User already exists\n"
	}
	err = json.NewEncoder(w).Encode(response) // json.encode writes the JSON encoding of response to w
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message+"sent UserID ", response["userId"], "to the client")
}
