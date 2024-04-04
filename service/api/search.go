package api

import (
	"errors"
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/Alessio-Olivieri/wasaProject/service/database"
	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

func (rt *_router) search_users(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "search_users:\ncallerId: " + strconv.FormatUint(ctx.UserId, 10) + "\n"

	// Get the user to search from the URL
	searched_user := r.URL.Query().Get("username")
	if searched_user == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(message + ErrUserNotExists.Error())
		return
	}
	message = message + "User to search: " + searched_user + "\n"

	// Get the usernames of the users that match the search
	users, err := rt.db.GetSearchUser(searched_user, ctx.UserId)
	if errors.Is(err, database.ErrUserNotExists) {
		ctx.Logger.Error(message + "Error: User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(message+"Error getting search user: ", err)
		return
	}

	users_list := ""

	for _, username := range users {
		users_list += username + ", "
	}
	message = message + "Users found: " + users_list + "\n"

	var search_result schemas.UserList
	search_result.Users = users
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(search_result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(message+"Error encoding search result: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info(message)
}
