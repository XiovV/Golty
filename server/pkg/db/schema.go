package db

const schema = `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  user_id INTEGER,
  refresh_token TEXT NOT NULL,
  FOREIGN KEY (user_id)
    REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS config (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
`

func (db *DB) initSchema() {
	db.db.MustExec(schema)
}
