package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetIdFromBearer(r *http.Request) uint64 {
	authHeader := r.Header.Get("authToken")
	fmt.Println("authheader:", authHeader)
	authParts := strings.Fields(authHeader)
	token := authParts[1]
	myId, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 18446744073709551615
	}
	return myId
}
