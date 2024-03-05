package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"

	"fmt"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getContextReply(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))

	//get cookie
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		ctx.Logger.Debug("no session cookie found")
	} else {
		ctx.SessionID = sessionCookie.Value
	}

	//parse sessionID
	fmt.Println("\n\n sessionID: \n\n", ctx.SessionID)
	sessionID, err := strconv.Atoi(ctx.SessionID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error parsing session id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get user from session
	userID, err := rt.db.Get_userID(sessionID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting user from session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get username from userID
	user, err := rt.db.Get_username(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting username from userID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(user))
}
