package requests

import "regexp"

var usernameSearchRegex = regexp.MustCompile(`^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*$`)

func check_usernameSearch(username string) bool {
	return 1 <= len(username) && len(username) <= 25 && usernameSearchRegex.MatchString(username)
}

type UsernameSearch struct {
	Username string `json:"username"`
}

func (request *UsernameSearch) IsValid() bool {
	return check_usernameSearch(request.Username)
}
