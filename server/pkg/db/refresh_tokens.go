package db

type RefreshToken struct {
	UserID       int    `db:"user_id"`
	RefreshToken string `db:"refresh_token"`
}

func (db *DB) InsertRefreshToken(userId int, refreshToken string) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := db.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, refresh_token) VALUES ($1, $2)", userId, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRefreshToken(userId int, refreshToken string) (string, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var foundToken RefreshToken
	err := db.db.GetContext(ctx, &foundToken, "SELECT user_id, refresh_token FROM refresh_tokens WHERE user_id = $1 AND refresh_token = $2", userId, refreshToken)
	if err != nil {
		return "", err
	}

	return foundToken.RefreshToken, nil
}

func (db *DB) UpdateRefreshToken(userId int, oldRefreshToken string, newRefreshToken string) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := db.db.ExecContext(ctx, "UPDATE refresh_tokens SET refresh_token = $1 WHERE user_id = $2 AND refresh_token = $3", newRefreshToken, userId, oldRefreshToken)
	if err != nil {
		return err
	}

	return nil
}
