package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"
)

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) put_like(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "put_like:\n"

	// Extract the photo_id from the path parameters
	photo_id_s := ps.ByName("photo_id")
	if photo_id_s == "" {
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "photo_id: " + photo_id_s + "\n"

	photo_id, err := strconv.ParseUint(photo_id_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error parsing photo_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// putting like
	err = rt.db.PutLike(photo_id, ctx.UserId)
	if errors.Is(err, database.ErrPhotoNotExists) {
		ctx.Logger.WithError(err).Error(message + "Error: photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrBanned) {
		ctx.Logger.WithError(err).Error(message + "Error: user is banned")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrAlreadyLiked) {
		ctx.Logger.Info(message + "OK: photo already liked")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting like in the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message + "Success")
}

// func (rt *_router) get_likes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var message string
// 	message = message + "get_likes:\n"

// 	// Extract the photo_id from the path parameters
// 	photo_id, err := strconv.ParseUint(ps.ByName("photo_id"), 10, 64)
// 	if err != nil {
// 		// Handle missing photo_id parameter
// 		ctx.Logger.Error(message + "Error: photo_id parameter missing in request path")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	message = message + "photo_id: " + strconv.FormatUint(photo_id, 10) + "\n"

// 	// Check if the photo exists
// 	exists, err := rt.db.Exists_photo(photo_id)
// 	if err != nil {
// 		// Handle database error
// 		ctx.Logger.WithError(err).Error(message + "Error checking if photo exists")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	if !exists {
// 		// Handle non-existent photo
// 		ctx.Logger.Error(message + "Error: photo does not exist")
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	ctx.Logger.Info(message + "Success")
// }

// UserId of follower is authenticated user,
// FollowedId is the id corresponding to the username of the visited profile
func (rt *_router) delete_like(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "delete_like:\n"

	// Extract the photo_id from the path parameters
	photo_id_s := ps.ByName("photo_id")
	if photo_id_s == "" {
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "photo_id: " + photo_id_s + "\n"

	photo_id, err := strconv.ParseUint(photo_id_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error parsing photo_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// deleting like
	err = rt.db.DeleteLike(photo_id, ctx.UserId)
	if errors.Is(err, database.ErrPhotoNotExists) {
		ctx.Logger.WithError(err).Error(message + "Error: photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrBanned) {
		ctx.Logger.WithError(err).Error(message + "Error: user is banned")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrLikeNotExists) {
		ctx.Logger.Info(message + "Like did not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting like in the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message)
}
