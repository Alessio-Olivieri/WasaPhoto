package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/components/requests"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "changeUsername:\n"
	var request requests.Username

	// decode json
	err := json.NewDecoder(r.Body).Decode(&request)
	_ = r.Body.Close()

	// check JSON err
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check correctness of input
	if !request.IsValid() {
		ctx.Logger.WithField("username", request.Username).Error(message + "username not valid - JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.Update_username(ctx.UserId, request.Username)
	if errors.Is(err, database.ErrUsernameTaken) {
		ctx.Logger.WithError(err).Error(message + "DATABASE Update_username:")
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "DATABASE Update_username:")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message)
}
