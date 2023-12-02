package main

import (
	"database/sql"
	"golty/repository"

	"github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

const (
	DefaultUsername = "admin"
	DefaultPassword = "admin"
	DefaultIsAdmin  = 1
)

var DefaultUser = repository.User{
	Username: DefaultUsername,
	Password: DefaultPassword,
	Admin:    DefaultIsAdmin,
}

var argon2Params = argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  10,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func initServer(repository *repository.Repository, logger *zap.Logger) error {
	logger.Debug("checking if any users exist...")
	_, err := repository.FindUserByUsername(DefaultUsername)
	if err != sql.ErrNoRows {
		logger.Debug("at least one user exists, continuing...")
		return nil
	}

	logger.Info("creating fresh admin user...")

	hash, err := argon2id.CreateHash(DefaultPassword, &argon2Params)
	if err != nil {
		return err
	}

	DefaultUser.Password = hash

	err = repository.InsertUser(DefaultUser)
	if err != nil {
		return err
	}

	logger.Info("fresh admin user created successfully")

	return nil
}
