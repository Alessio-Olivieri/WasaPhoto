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
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	AddUser(name string) error
	ListUsers() ([]string, error)
	Exists_user(username string) (int, error)

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
