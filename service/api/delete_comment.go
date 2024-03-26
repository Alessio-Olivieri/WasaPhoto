package api

import (
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) delete_comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "post_photo: authUser:" + strconv.FormatUint(ctx.UserId, 10) + "\n"

	comment_id_s := ps.ByName("comment_id")
	if comment_id_s == "" {
		ctx.Logger.Error(message + "comment_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "comment_id_s:" + comment_id_s + "\n"

	comment_id, err := strconv.ParseUint(comment_id_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing comment_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(comment_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error posting comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message)
}
