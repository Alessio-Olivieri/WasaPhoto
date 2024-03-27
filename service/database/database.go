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
	"fmt"
	"log"
	"mime/multipart"

	"github.com/Alessio-Olivieri/wasaProject/service/schemas"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Get the username of the user with the specified id
	Get_userId_from_username(user_id string) (uint64, error)
	Get_username_from_userId(userId uint64) (string, error)
	ListUsers() ([]string, error)
	GetSearchUser(prompt string, requestingUser uint64) ([]string, error)
	//if user exists return its ID, otherwise create a new user and return its ID
	Exists_user(username string) (int, error)

	// Bans
	//if banner has banned banned return 1, otherwise 0
	IsBanned(banner uint64, banned uint64) (bool, error)
	// insert a ban in the database
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
	//determine if follower is following followed
	IsFollowing(follower_id uint64, followed_id uint64) (bool, error)
	GetFollowersAmount(followed_id uint64) (int, error)

	// Photo
	// Make_photo creates a new photo in the database
	Make_photo(user_id uint64, caption string, picture multipart.File) (schemas.Post, error)
	Get_user_photos(user_id_request uint64, user_id_profile uint64) ([]schemas.Post, error)
	Exists_photo(photo_id uint64) (bool, error)

	// likes
	DeleteLike(photo_id uint64, user_id uint64) error
	PutLike(photo_id uint64, user_id uint64) error
	IsLiked(photo_id uint64, user_id uint64) (bool, error)
	//get likes
	GetLikeAmount(photo_id uint64) (int, error)

	// comments
	PostComment(photo_id uint64, user_id uint64, text string) error
	DeleteComment(comment_id uint64) error
	GetComments(photo_id uint64) ([]schemas.Comment, error)

	// stream
	Get_stream(request_user_id uint64, page int) ([]schemas.Post, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	log.Println("[BUILDING DATABASE STRUCTURE]")
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
		fmt.Printf(" checking for table %s:\n", tableName)
		// Check if table exists. If not, the database is empty, and we need to create the structure

		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name= ? ;`, tableName).Scan(&tableName)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("  creating database %s\n", tableName)
			_, err = db.Exec(tableCreationStatement)
			if err != nil {
				return nil, fmt.Errorf("   error creating database structure: %w", err)
			}
		} else {
			fmt.Printf("  table %s already exists\n", tableName)
		}
	}
	log.Println("[DATABASE STRUCTURE CREATED]")

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
