package api

import (
	"errors"
	"strconv"
)

func ExtractId_from_Bearer(token string) (uint64, error) {
	if len(token) < len("Bearer ") || token[:len("Bearer ")] != "Bearer " {
		return 0, errors.New("invalid Bearer token format")
	}
	return strconv.ParseUint(token[len("Bearer "):], 10, 64)
}
