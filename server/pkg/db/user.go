package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (db *DB) initUser() {
	var user User

	err := db.db.Get(&user, "SELECT email, password FROM users LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}

		db.logger.Fatal("an unknown error occured while checking for existing users", "error", err)
		return
	}
}

func (db *DB) createDefaultUser() {
	_, err := db.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "admin@email.com", "test123")
	if err != nil {
		fmt.Println("error", err)
	}
}
