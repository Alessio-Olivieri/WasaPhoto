package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) Get_username_from_userId(userId uint64) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM Users WHERE user_id = ?", userId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
