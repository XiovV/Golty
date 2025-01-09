package db

import (
	"database/sql"

	"github.com/matthewhartstonge/argon2"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password []byte `db:"password"`
}

func (db *DB) InitUser() error {
	var existingUser User

	err := db.db.Get(&existingUser, "SELECT username, password FROM users LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows {
			db.logger.Info("no users found, generating default admin user")
			return db.createDefaultUser()
		}

		db.logger.Fatal("an unknown error occured while checking for existing users", "error", err)
		return err
	}

	return nil
}

func (db *DB) createDefaultUser() error {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(DEFAULT_ADMIN_PASSWORD))
	if err != nil {
		return err
	}

	_, err = db.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", DEFAULT_ADMIN_USERNAME, encoded)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserByUsername(username string) (User, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var user User

	err := db.db.GetContext(ctx, &user, "SELECT id, username, password FROM users WHERE username = $1", username)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
