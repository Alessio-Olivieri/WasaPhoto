package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type objectID struct {
	OID uint64 `json:"$OID"`
}

func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	ctx.Logger.Info("Get profile")

	// Extract the username from the path parameters using httprouter.Params
	username := ps.ByName("username")
	if username == "" {
		// Handle missing username parameter
		ctx.Logger.Error("Error: Username parameter missing in request path")
		return
	}

	ctx.Logger.Info("Extracted username: ", username)

	// Check if the auth token matches the userID of username
	userID, err := rt.db.Get_userId(username)
	if err != nil || userID == 18446744073709551615 {
		ctx.Logger.WithError(err).Error("Error getting user ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Logger.Info("User ID: ", userID)

	// var authToken string
	// ctx.Logger.Info("authToken: ", authToken)
	// if err := json.NewDecoder(r.Body).Decode(&authToken); err != nil {
	// 	ctx.Logger.WithError(err).Error("Error decoding authtoken from request body")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// context_userID, err := strconv.ParseUint(strings.Fields(authToken)[1], 10, 64)
	// if err != nil {
	// 	ctx.Logger.WithError(err).Error("Error parsing user ID from request body")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// if userID != context_userID {
	// 	ctx.Logger.Error("Error: User ID does not match the auth token : ", userID, "!=", context_userID)
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	//get user photos

}
