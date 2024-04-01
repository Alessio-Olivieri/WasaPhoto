package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

// login is the handler for the POST /login route
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "login:\n"

	// request body
	var loginReq schemas.Username
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		ctx.Logger.WithError(err).Error(message + "Error parsing json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database, if it does not exist, create it
	userID, err := rt.db.Login_db(loginReq.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error(message + "error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Authenticated user " + loginReq.Username + " with ID " + strconv.FormatUint(userID, 10) + "\n"

	// response contains the user ID associated to the username
	response := map[string]uint64{"userId": userID}

	// Send the response to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response) // json.encode writes the JSON encoding of response to w
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message+"sent UserID ", response["userId"], "to the client")
}
