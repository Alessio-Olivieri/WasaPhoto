package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// LoginRequest is the request body for the POST /login route
type LoginRequest struct {
	Username string `json:"username"`
}

// login is the handler for the POST /login route
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Info("Login request")

	// Decode the request body
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database, if it does not exist, create it
	userID, err := rt.db.Exists_user(loginReq.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// response contains the user ID associated to the username
	response := map[string]uint64{"userId": uint64(userID)}

	// Send the response to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response) // json.encode writes the JSON encoding of response to w
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't login user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
