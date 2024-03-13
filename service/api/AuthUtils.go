package api

import (
	"net/http"
	"strconv"
	"strings"
)

func GetIdFromBearer(r *http.Request) uint64 {
	authHeader := r.Header.Get("Authorization")
	authParts := strings.Fields(authHeader)
	token := authParts[1]
	myId, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 100000
	}
	return myId
}
