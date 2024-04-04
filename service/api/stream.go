package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) get_stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var message string
	message = message + "getUserProfile:\n"

	// Get the page number from the URL
	page_number, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error parsing page number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if page_number < 0 {
		ctx.Logger.WithError(err).Error(message + "page number can't be >0")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message = message + "Page number: " + strconv.FormatInt(page_number, 10) + "\n"

	// retrieve the stream
	stream, err := rt.db.Get_stream(ctx.UserId, int(page_number))
	if errors.Is(err, database.ErrEmptyStream) {
		ctx.Logger.WithError(err).Error(message + "No posts in the stream")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, database.ErrEmptyStream) {
		ctx.Logger.WithError(err).Error(message + "Page out of bound")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error retrieving posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message = message + strconv.Itoa(len(stream.Posts)) + " posts retrieved successfully\n"

	stream.Page_number = int(page_number)

	// create json response
	jsonData, err := json.Marshal(stream)
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

	ctx.Logger.Info(message)
}
