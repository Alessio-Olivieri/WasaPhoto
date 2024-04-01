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

func (rt *_router) post_photo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "post_photo: authUser:" + strconv.FormatUint(ctx.UserId, 10) + "\n"

	// parse request
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.Error(message + "post-photo: error while parsing multipart form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get caption from request
	caption := r.FormValue("caption")
	emptyComment := regexp.MustCompile(`^\s*(?:#.*|\bundefined\b|\s*)$`)
	empty_caption := false
	if emptyComment.MatchString(caption) {
		message = message + "empty caption \n"
		empty_caption = true
	}

	//get photo from request
	file, _, err := r.FormFile("picture")
	empty_picture := false
	if err != nil {
		if errors.Is(http.ErrMissingFile, err) {
			message = message + "No photo present \n"
			empty_picture = true
		}
		if !errors.Is(http.ErrMissingFile, err) {
			ctx.Logger.Error(message + "error while obtaining photo resource " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}

	if empty_caption && empty_picture {
		ctx.Logger.Error(message + "No photo and caption present")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := rt.db.Make_photo(ctx.UserId, caption, file)
	if err != nil {
		ctx.Logger.Error(message + "Database error: Make_photo:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Photo inserted in database \n"

	// create json response
	jsonData, err := json.Marshal(response)
	if err != nil {
		ctx.Logger.Error(message + "Error while creating json response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Json response created \n"

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	amount, err := w.Write(jsonData)
	if err != nil {
		ctx.Logger.Error(message + "Error while writing response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Response written: " + strconv.Itoa(amount) + " bytes \n"

	w.WriteHeader(http.StatusCreated)
	ctx.Logger.Info(message)
}

func (rt *_router) delete_photo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "delete_photo:\n"

	// Extract the photo_id from the path parameters
	photo_id_s := ps.ByName("photo_id")
	if photo_id_s == "" {
		ctx.Logger.Error(message + "Error: Photo_id parameter missing in request path")
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

	//deleting photo
	err = rt.db.DeletePost(photo_id, ctx.UserId)
	if errors.Is(err, database.ErrPhotoNotExists) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeletePost:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeletePost:")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "DATABASE ERROR: DeletePost:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "Photo deleted successfully")
	w.WriteHeader(http.StatusOK)
}
