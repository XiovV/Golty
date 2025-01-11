package db

type Config struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

func (db *DB) InsertConfig(key string, value string) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := db.db.ExecContext(ctx, "INSERT INTO config (key, value) VALUES ($1, $2)", key, value)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetConfig(key string) (string, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var config Config
	err := db.db.GetContext(ctx, &config, "SELECT key, value FROM config WHERE key = $1 ", key)
	if err != nil {
		return "", err
	}

	return config.Value, nil
}
