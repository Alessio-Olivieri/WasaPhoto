package database

import (
	"database/sql"
	"errors"
)

// retruns 1 if banned, 0 if not banned, -1 if error
func (db *appdbimpl) IsLiked(photo_id uint64, user_id uint64) (bool, error) {
	var exists int8
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Likes
		WHERE photo_id = ? AND user_id = ?)`,
		photo_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func (db *appdbimpl) PutLike(photo_id uint64, user_id uint64) error {
	var user_id_photo uint64
	err := db.c.QueryRow(`SELECT user_id FROM Photos WHERE photo_id = ?`, photo_id).Scan(&user_id_photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPhotoNotExists
		}
		return err
	}

	// check if user is banned
	banned, err := db.IsBanned(user_id_photo, user_id)
	if err != nil {
		return err
	}
	if banned {
		return ErrBanned
	}

	// check if like exists
	exists, err := db.IsLiked(photo_id, user_id)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyLiked
	}

	_, err = db.c.Exec(`INSERT into Likes VALUES (?, ?)`, photo_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteLike(photo_id uint64, user_id uint64) error {
	var user_id_photo uint64
	err := db.c.QueryRow(`SELECT user_id FROM Photos WHERE photo_id = ?`, photo_id).Scan(&user_id_photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPhotoNotExists
		}
		return err
	}

	// check if user is banned
	banned, err := db.IsBanned(user_id_photo, user_id)
	if err != nil {
		return err
	}
	if banned {
		return ErrBanned
	}

	// check if like exists
	exists, err := db.IsLiked(photo_id, user_id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrLikeNotExists
	}
	_, err = db.c.Exec(`DELETE from Likes WHERE photo_id = ? AND user_id = ?`, photo_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

// get like amount
func (db *appdbimpl) GetLikeAmount(photo_id uint64) (int, error) {
	var like_count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM Likes
		WHERE photo_id = ?`, photo_id).Scan(&like_count)
	if err != nil {
		return 0, err
	}
	return like_count, nil
}

// get the usernames of who putted like
func (db *appdbimpl) GetLikes(photo_id uint64) ([]string, error) {
	var likes []string
	rows, err := db.c.Query(`SELECT user_id FROM Likes
		WHERE photo_id = ?`, photo_id)
	if err != nil {
		return likes, err
	}
	defer rows.Close()

	for rows.Next() {
		var likeId uint64
		err = rows.Scan(&likeId)
		if err != nil {
			return nil, err
		}

		like, err := db.Get_username_from_userId(likeId)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	return likes, nil
}
