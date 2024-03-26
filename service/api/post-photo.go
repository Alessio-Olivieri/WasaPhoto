package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
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
	file, _, err := r.FormFile("picture")
	if err != nil {
		if err.Error() == "http: no such file" && caption == "" {
			ctx.Logger.Error(message + "no file and no caption")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err.Error() == "http: no such file" {
			message = message + "No photo present \n"
		}
		if err.Error() != "http: no such file" {
			ctx.Logger.Error(message + "error while obtaining photo resource " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if file != nil {
		defer file.Close()
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

	ctx.Logger.Info(message)
}
