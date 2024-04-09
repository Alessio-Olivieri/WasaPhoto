package api

import (
	"errors"
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"
)

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) put_follower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "put_Follower:\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		return
	}
	message = message + profile_username + "\n"

	profileId, err := rt.db.Get_userId_from_username(profile_username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(message + "Error getting userId from Username")
		return
	}

	if profileId == ctx.UserId {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error(message + "you can't follow yourself")
		return
	}

	// check if you banned the user
	banned, err := rt.db.IsBanned(ctx.UserId, profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error checking if user is banned")
		return
	}
	if banned {
		ctx.Logger.Error(message + "User has banned been banned, before following unban him first")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user banned you
	banned, err = rt.db.IsBanned(profileId, ctx.UserId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error checking if user is banned")
		return
	}
	if banned {
		ctx.Logger.Error(message + "User has banned you")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// putting follow
	err = rt.db.PutFollower(ctx.UserId, profileId)
	if errors.Is(database.ErrAlreadyFollowed, err) {
		ctx.Logger.Info(message + "OK: User is already following the profile")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting follower in the database")
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message)
}

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) delete_follower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "Delete_follower" + "\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		return
	}
	message = message + profile_username + "\n"

	profileId, err := rt.db.Get_userId_from_username(profile_username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(message + "Error getting userId from Username")
		return
	}

	// check if you banned the user
	banned, err := rt.db.IsBanned(ctx.UserId, profileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(message + "Error checking if user is banned")
		return
	}
	if banned {
		ctx.Logger.Error(message + "User has banned been banned, before following unban him first")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user banned you
	banned, err = rt.db.IsBanned(profileId, ctx.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(message + "Error checking if user is banned")
		return
	}
	if banned {
		ctx.Logger.Error(message + "User has banned you")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// deleteing follow
	err = rt.db.DeleteFollower(ctx.UserId, profileId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error deleting follower from the database")
		return
	}
	if errors.Is(database.ErrFollowNotExists, err) {
		ctx.Logger.Info(message + "OK: User is already not following the profile")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message)
}
