package db

func (db *DB) InsertRefreshToken(userId int, refreshToken string) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := db.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, refresh_token) VALUES ($1, $2)", userId, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
