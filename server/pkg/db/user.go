package db

import (
	"database/sql"

	"github.com/matthewhartstonge/argon2"
)

type User struct {
	Email    string `db:"email"`
	Password []byte `db:"password"`
}

func (db *DB) InitUser() error {
	var existingUser User

	err := db.db.Get(&existingUser, "SELECT email, password FROM users LIMIT 1")
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
