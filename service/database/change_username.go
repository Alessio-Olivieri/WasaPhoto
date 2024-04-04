package database

func (db *appdbimpl) Update_username(userId uint64, username string) error {
	// check if username is already taken
	var exists bool
	exists, err := db.Exists_user(username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUsernameTaken
	}

	_, err = db.c.Exec("Update Users SET username = ? WHERE user_id = ?", username, userId)
	if err != nil {
		return err
	}
	return nil

}
