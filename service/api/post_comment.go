package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type CommentReq struct {
	Text string `json:"content"`
}

func (rt *_router) post_comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "post_photo: authUser:" + strconv.FormatUint(ctx.UserId, 10) + "\n"

	photoId_s := ps.ByName("photo_id")
	if photoId_s == "" {
		ctx.Logger.Error(message + "photo_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "photo_id:" + photoId_s + "\n"

	photoId, err := strconv.ParseUint(photoId_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing photo_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var commentReq CommentReq
	if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "text:" + commentReq.Text + "\n"

	err = rt.db.PostComment(photoId, ctx.UserId, commentReq.Text)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error posting comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message)
}
