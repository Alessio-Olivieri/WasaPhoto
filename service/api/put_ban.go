package api

import (
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) put_ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "put_Ban:\n"

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
	message = message + "profile_username: " + profile_username + "\n"

	//putting ban
	err = rt.db.PutBan(ctx.UserId, profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting ban in the database")
		return
	}

	//removing follow
	err = rt.db.DeleteFollower(ctx.UserId, profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error deleting follow from the database")
		return
	}

	ctx.Logger.Info(message)
}
