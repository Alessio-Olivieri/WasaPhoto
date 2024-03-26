package api

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
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
	message = message + "Page number: " + strconv.FormatInt(page_number, 10) + "\n"

	//retrieve the stream
	response, err := rt.db.Get_stream(ctx.UserId, int(page_number))
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "Error retrieving posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + strconv.Itoa(len(response)) + " posts retrieved successfully\n"

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

	ctx.Logger.Info(message)
}
