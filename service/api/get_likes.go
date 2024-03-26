package api

import (
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) get_likes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "get_likes:\n"

	// Extract the photo_id from the path parameters
	photo_id, err := strconv.ParseUint(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		// Handle missing photo_id parameter
		ctx.Logger.Error(message + "Error: photo_id parameter missing in request path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "photo_id: " + strconv.FormatUint(photo_id, 10) + "\n"

	// Check if the photo exists
	exists, err := rt.db.Exists_photo(photo_id)
	if err != nil {
		// Handle database error
		ctx.Logger.WithError(err).Error(message + "Error checking if photo exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		// Handle non-existent photo
		ctx.Logger.Error(message + "Error: photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the list of users who liked the photo
	likes, err := rt.db.GetLikeAmount(photo_id)
	if err != nil {
		// Handle database error
		ctx.Logger.WithError(err).Error(message + "Error getting likes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the list of likes
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(strconv.Itoa(likes)))
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info(message + "Success")
}
