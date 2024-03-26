package api

import (
	"net/http"

	"github.com/Alessio-Olivieri/wasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"strconv"

	"github.com/Alessio-Olivieri/wasaProject/service/schemas"
)

func (rt *_router) search_users(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var message string
	message = message + "Called search_users, callerId: " + strconv.FormatUint(ctx.UserId, 10) + "\n"

	// Get the user to search from the URL
	searched_user := r.URL.Query().Get("username")
	message = message + "User to search: " + searched_user + "\n"

	// Get the usernames of the users that match the search
	users, err := rt.db.GetSearchUser(searched_user, ctx.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error(message+"Error getting search user: ", err)
		return
	}

	users_list := ""

	for _, username := range users {
		//fmt.Println(username) //display list in terminal
		users_list += username + ", "
	}
	message = message + "Users found: " + users_list + "\n"

	var search_result schemas.UserList
	search_result.Users = users
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(search_result)

	ctx.Logger.Info(message)
}
