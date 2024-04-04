package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
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
	emptyComment := regexp.MustCompile(`^\s*(?:#.*|\bundefined\b|\s*)$`)
	if emptyComment.MatchString(commentReq.Text) {
		ctx.Logger.Error(message + "empty comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "text:" + commentReq.Text + "\n"

	comment, err := rt.db.PostComment(photoId, ctx.UserId, commentReq.Text)
	if errors.Is(err, database.ErrPhotoNotExists) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: PostComment:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrBanned) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: PostComment:")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error posting comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	ctx.Logger.Info(message)
}

// func (rt *_router) get_comments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

// }

func (rt *_router) delete_comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "delete comment: authUser:" + strconv.FormatUint(ctx.UserId, 10) + "\n"

	comment_id_s := ps.ByName("comment_id")
	if comment_id_s == "" {
		ctx.Logger.Error(message + "comment_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment_id, err := strconv.ParseUint(comment_id_s, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing comment_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "comment_id:" + comment_id_s + "\n"

	err = rt.db.DeleteComment(comment_id, ctx.UserId)
	if errors.Is(err, database.ErrCommentNotExists) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeleteComment:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeleteComment:")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeleteComment:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message + "Comment deleted successfully")
}
