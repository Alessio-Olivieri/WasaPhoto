package database

// retruns 1 if banned, 0 if not banned, -1 if error
func (db *appdbimpl) IsBanned(banner_id uint64, banned_id uint64) (bool, error) {
	var exists int8
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Bans
		WHERE banner_id = ? AND banned_id = ?)`,
		banner_id, banned_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

// inserts a new ban in the database if it does not already exist
func (db *appdbimpl) PutBan(banner_id uint64, banned_id uint64) error {
	// Check if ban exists
	exists, err := db.IsBanned(banner_id, banned_id)
	if err != nil {
		return err
	}
	if exists {
		return ErrBanAlreadyExists
	}
	_, err = db.c.Exec(`INSERT INTO Bans (banner_id, banned_id) VALUES (?, ?);`, banner_id, banned_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteBan(banner_id uint64, banned_id uint64) error {
	// Check if ban exists
	exists, err := db.IsBanned(banner_id, banned_id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrBanNotExists
	}

	_, err = db.c.Exec(`DELETE from Bans WHERE banner_id = ? AND banned_id = ?`, banner_id, banned_id)
	if err != nil {
		return err
	}
	return nil
}
