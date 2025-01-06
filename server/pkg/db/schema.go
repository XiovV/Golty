package db

const schema = `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password BLOB NOT NULL
)
`

func (db *DB) initSchema() {
	db.db.MustExec(schema)
}
