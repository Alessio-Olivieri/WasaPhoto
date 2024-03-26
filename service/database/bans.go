package database

// retruns 1 if banned, 0 if not banned, -1 if error
func (db *appdbimpl) IsBanned(bannerId uint64, bannedId uint64) (bool, error) {
	var exists int8
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Bans
		WHERE banner_id = ? AND banned_id = ?)`,
		bannerId, bannedId).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func (db *appdbimpl) PutBan(follower_id uint64, followed_id uint64) error {
	_, err := db.c.Exec(`INSERT into Bans VALUES (?, ?)`, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteBan(follower_id uint64, followed_id uint64) error {
	_, err := db.c.Exec(`DELETE from Bans WHERE banner_id = ? AND banned_id = ?`, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}
