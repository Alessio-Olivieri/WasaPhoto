package api

import (
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) delete_follower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "Delete_follower" + "\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		ctx.Logger.Error(message + "Error: Username parameter missing in request path")
		return
	}
	message = message + "profile_username: " + profile_username + "\n"

	profileId, err := rt.db.Get_userId_from_username(profile_username)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting userId from Username")
		return
	}

	//putting follow
	err = rt.db.DeleteFollower(ctx.UserId, profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error deleting follower from the database")
		return
	}

	//return the number of followers
	followers, err := rt.db.GetFollowersAmount(profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting followers")
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(strconv.Itoa(followers)))
	w.WriteHeader(http.StatusOK)

	ctx.Logger.Info(message)
}
