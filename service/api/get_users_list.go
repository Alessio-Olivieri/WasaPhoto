package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) get_users_list(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("users list: "))
	log.Println("get_users_list called!")

	users, err := rt.db.ListUsers()
	if err != nil {
		log.Println(" error listing users ", err)
	}

	users_list := ""

	for _, username := range users {
		//fmt.Println(username) //display list in terminal
		users_list += username + ", "
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte(users_list))

}
