package database

import (
	"errors"
	"strconv"
)

func (db *appdbimpl) Get_user_photos(user_id uint64) (uint64, error) {
	user_id_s := strconv.FormatUint(user_id, 64)
	if user_id_s == "" {
		return 0, errors.New("user_id is blank")
	}
	return 0, nil
}
