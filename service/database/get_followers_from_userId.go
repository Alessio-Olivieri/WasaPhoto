package database

import (
	"strconv"
)

func (db *appdbimpl) Get_followers_from_userId(userId uint64) ([]string, error) {
	var followers []string
	rows, err := db.c.Query("SELECT follower_id FROM Followers WHERE followed_id = ?", userId)
	if err != nil {
		return followers, err
	}
	defer rows.Close()

	// Iterate over the rows, getting the followers
	for rows.Next() {
		var followerId string
		err := rows.Scan(&followerId)
		if err != nil {
			return followers, err
		}
		if rows.Err() != nil {
			return followers, rows.Err()
		}
		// Get the username of the follower
		followerId_int, err := strconv.ParseUint(followerId, 10, 64)
		if err != nil {
			return followers, err
		}
		follower_name, err := db.Get_username_from_userId(followerId_int)
		if err != nil {
			return followers, err
		}
		// Append the username to the followers list
		followers = append(followers, follower_name)
	}
	// Return the list of followers
	return followers, nil
}
