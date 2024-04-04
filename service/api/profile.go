package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var message string
	message = message + "getUserProfile:\n"

	// Extract the username from the path parameters
	profile_username := ps.ByName("username")
	if profile_username == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(message + "Error: Username parameter missing in request path")
		return
	}
	message = message + "profile_username: " + profile_username + "\n"

	// retrieve the userId from the username in the path
	userId_profile, err := rt.db.Get_userId_from_username(profile_username)
	if errors.Is(err, database.ErrUserNotExists) {
		ctx.Logger.Error(message + "Error: User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil || userId_profile == 18446744073709551615 {
		ctx.Logger.WithError(err).Error(message + "Error getting userId from profile_username ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + " UserId_profile: " + strconv.FormatUint(userId_profile, 10) + "\n"

	// Check if the user has banned you
	banned, err := rt.db.IsBanned(userId_profile, ctx.UserId)
	if err != nil {
		ctx.Logger.Error(message + "" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		ctx.Logger.Error(message + "Error: User_profile has banned you")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	message = message + "User has not banned you" + "\n"

	// build the response
	var response schemas.GetUserProfile

	// Get profile followers
	followers, err := rt.db.Get_followers_from_userId(userId_profile)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting followers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Followers = followers
	message = message + "Followers: " + strconv.Itoa(len(followers)) + "\n"

	// add followerscount to response
	response.Followers_count = len(followers)

	// determine if the reqeusting user is following the profile_user
	isFollowing, err := rt.db.IsFollowing(ctx.UserId, userId_profile)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error determining follower status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.IsFollowing = isFollowing
	message = message + "Requesting user is not following " + profile_username + "\n"

	// determine if the requesting user has banned the profile_user
	isBanned, err := rt.db.IsBanned(ctx.UserId, userId_profile)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error determining ban status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.IsBanned = isBanned
	message = message + "Requesting user has not banned " + profile_username + "\n"

	// retrieve photos of the user_profile
	post_list, err := rt.db.Get_user_photos(ctx.UserId, userId_profile)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error retrieving posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Posts = post_list
	message = message + strconv.Itoa(len(post_list)) + " photos of " + profile_username + " retrieved successfully\n"

	// create json response
	jsonData, err := json.Marshal(response)
	if err != nil {
		ctx.Logger.Error(message + "Error while creating json response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Json response created \n"

	w.Header().Set("Content-Type", "application/json")
	amount, err := w.Write(jsonData)
	if err != nil {
		ctx.Logger.Error(message + "Error while writing response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Response written: " + strconv.Itoa(amount) + " bytes \n"

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info(message)
}
