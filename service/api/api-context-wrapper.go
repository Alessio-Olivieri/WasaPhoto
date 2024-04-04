package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func ExtractId_from_Bearer(token string) (uint64, error) {
	if len(token) < len("Bearer ") || token[:len("Bearer ")] != "Bearer " {
		return 0, errors.New("invalid Bearer token format")
	}
	return strconv.ParseUint(token[len("Bearer "):], 10, 64)
}

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler, authRequired bool) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Add the user id to the context if it's available
		if authRequired {
			ctx.UserId, err = ExtractId_from_Bearer(r.Header.Get("Authorization"))
			if err != nil {
				ctx.Logger.WithError(err).Error("ERROR wrap: ")
				ctx.UserId = 0
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
