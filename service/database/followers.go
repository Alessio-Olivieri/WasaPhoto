package database

func (db *appdbimpl) PutFollower(follower_id uint64, followed_id uint64) error {
	//check if follow exists
	exists, err := db.IsFollowing(follower_id, followed_id)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyFollowed
	}

	_, err = db.c.Exec(`INSERT into Followers VALUES (?, ?)`, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) IsFollowing(follower_id uint64, followed_id uint64) (bool, error) {
	rows, err := db.c.Query(`SELECT EXISTS(SELECT 1 FROM Followers
		WHERE follower_id = ? AND followed_id = ?)`,
		follower_id, followed_id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var result bool

	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (db *appdbimpl) GetFollowersAmount(followed_id uint64) (int, error) {
	rows, err := db.c.Query(`SELECT COUNT(*) FROM Followers WHERE followed_id = ?`, followed_id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var result int

	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (db *appdbimpl) DeleteFollower(follower_id uint64, followed_id uint64) error {
	//check if follow exists
	exists, err := db.IsFollowing(follower_id, followed_id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrFollowNotExists
	}

	_, err = db.c.Exec(`DELETE from Followers WHERE  follower_id = ? AND followed_id = ?`, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}
