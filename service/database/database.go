/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"mime/multipart"

	"github.com/Alessio-Olivieri/wasaProject/service/components/schemas"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Users

	// Get the id of the user with the specified username
	// returns ErrUserNotExists if the user doesn't exist
	Get_userId_from_username(user_id string) (uint64, error)
	// Get the username of the user with the specified id
	Get_username_from_userId(userId uint64) (string, error)
	// Get the list of all users
	ListUsers() ([]string, error)
	// Get the list of all users that match the prompt
	GetSearchUser(prompt string, requestingUser uint64) ([]string, error)
	// if user exists return its ID, otherwise create a new user and return its ID
	Login_db(username string) (uint64, error)
	// if user exists return true, otherwise return false
	Exists_user(identifier interface{}) (bool, error)
	// update the username of the user with the specified id
	Update_username(userId uint64, username string) error

	// Bans
	// if banner has banned the user specified by banned return true, otherwise false
	IsBanned(banner uint64, banned uint64) (bool, error)
	// inserts a new ban in the database if it does not already exist
	PutBan(follower_id uint64, followed_id uint64) error
	// delete a ban from the database
	DeleteBan(follower_id uint64, followed_id uint64) error

	// Followers
	// returns the list of followers of the user with the specified id
	Get_followers_from_userId(userId uint64) ([]string, error)
	// insert a follow in the database
	PutFollower(follower_id uint64, followed_id uint64) error
	// delete a follow from the database
	DeleteFollower(follower_id uint64, followed_id uint64) error
	// determine if follower is following followed
	IsFollowing(follower_id uint64, followed_id uint64) (bool, error)
	// get the number of followers of the user with the specified id
	GetFollowersAmount(followed_id uint64) (int, error)

	// Photo
	// Make_photo creates a new photo in the database
	Make_photo(user_id uint64, caption string, picture multipart.File) (schemas.Post, error)
	// Get_photo returns the photos of the user with the specified id, if the user requesting is banned return an error
	Get_user_photos(user_id_request uint64, user_id_profile uint64) ([]schemas.Post, error)
	// exists_photo returns true if the photo with the specified id exists, otherwise false
	Exists_photo(photo_id uint64) (bool, error)
	// delete_photo deletes the photo with the specified id, if the user requesting is not owner of photo return an error
	DeletePost(post_id uint64, user_id uint64) error

	// Likes
	// delete a like from the database
	DeleteLike(photo_id uint64, user_id uint64) error
	// insert a like in the database
	PutLike(photo_id uint64, user_id uint64) error
	// determine if user has liked photo
	IsLiked(photo_id uint64, user_id uint64) (bool, error)
	// get the number of likes of the photo with the specified id
	GetLikeAmount(photo_id uint64) (int, error)

	// comments
	// insert a comment in the database if the user is not banned
	PostComment(photo_id uint64, user_id uint64, text string) (schemas.Comment, error)
	// delete a comment from the database
	DeleteComment(comment_id uint64, user_id uint64) error
	// get the comments of the photo with the specified id
	GetComments(photo_id uint64) ([]schemas.Comment, error)

	// stream
	Get_stream(request_user_id uint64, page int) (schemas.Stream, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	// message := "[BUILDING DATABASE STRUCTURE]\n"
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	TableMapping := map[string]string{
		"Users":     usersTableCreationStatement,
		"Photos":    photosTableCreationStatement,
		"Comments":  commentsTableCreationStatement,
		"Likes":     likesTableCreationStatement,
		"Followers": followersTableCreationStatement,
		"Bans":      bansTableCreationStatement,
	}
	for tableName, tableCreationStatement := range TableMapping {
		// Check if table exists. If not, we need to create the structure
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name= ? ;`, tableName).Scan(&tableName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// message = message + "creating table: " + tableName + "\n"
				_, err = db.Exec(tableCreationStatement)
				if err != nil {
					return nil, errors.New("error building table " + tableName)
				}
			} else {
				return nil, errors.New("error checking table " + tableName)
			}
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
