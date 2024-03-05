package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) Get_username(user_id int) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT username FROM Users WHERE user_id = ?", user_id).Scan(&name)
	return name, err
}
