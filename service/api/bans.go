package api

import (
	"errors"
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) put_ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "put_Ban:\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileId, err := rt.db.Get_userId_from_username(profile_username)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting userId from Username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + profile_username + "\n"

	// putting ban
	err = rt.db.PutBan(ctx.UserId, profileId)
	if errors.Is(err, database.ErrBanAlreadyExists) {
		ctx.Logger.WithError(err).Error(message + "User was already banned\n")
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting ban in the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// removing follow
	err = rt.db.DeleteFollower(ctx.UserId, profileId)
	if errors.Is(err, database.ErrFollowNotExists) {
		message = message + "User was already unfollowed\n"
	} else if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error deleting follow from the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	ctx.Logger.Info(message + "Ban inserted in the database")
}

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) delete_ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "Delete_Ban" + "\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + profile_username + "\n"

	profileId, err := rt.db.Get_userId_from_username(profile_username)
	if errors.Is(err, database.ErrUserNotExists) {
		ctx.Logger.WithError(err).Error(message)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting userId from Username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete ban
	err = rt.db.DeleteBan(ctx.UserId, profileId)
	if errors.Is(err, database.ErrBanNotExists) {
		ctx.Logger.WithError(err).Error("DATABASE DeleteBan")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error deleting ban from the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info(message)
}
