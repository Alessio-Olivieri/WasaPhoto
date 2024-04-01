package requests

import (
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*[\w]$`)

type Username struct {
	Username string `json:"username"`
}

func (request *Username) IsValid() bool {
	return 3 <= len(request.Username) && len(request.Username) <= 25 && usernameRegex.MatchString(request.Username)
}
