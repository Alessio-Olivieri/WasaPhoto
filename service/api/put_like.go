package api

import (
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
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
		ctx.Logger.Error(message + "Error: Username parameter missing in request path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "photo_id: " + photo_id_s + "\n"

	photo_id, err := strconv.ParseUint(photo_id_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error parsing photo_id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//putting like
	err = rt.db.PutLike(photo_id, ctx.UserId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error inserting like in the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//return the number of likes
	likes, err := rt.db.GetLikeAmount(photo_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error getting likes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(strconv.Itoa(likes)))
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info(message)
}
