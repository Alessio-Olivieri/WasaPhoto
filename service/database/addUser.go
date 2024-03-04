package database

// function to populate the database of users with some initial data
func (db *appdbimpl) AddUser(username string) error {
	_, err := db.c.Exec(`INSERT INTO users (username) VALUES (?);`, username)
	return err
}
