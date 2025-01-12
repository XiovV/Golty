package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"github.com/XiovV/Golty/pkg/db"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/matthewhartstonge/argon2"
)

func (s *Server) Bootstrap() error {
	err := s.initializeDefaultUser()
	if err != nil {
		return err
	}

	secret, err := s.initializeTokenSecret()
	if err != nil {
		return err
	}

	s.jwtConfig = echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}

	return nil
}

func (s *Server) initializeDefaultUser() error {
	count, err := s.db.GetNumberOfUsers()
	if err != nil {
		return err
	}

	if count >= 1 {
		return nil
	}

	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(DEFAULT_ADMIN_PASSWORD))
	if err != nil {
		return err
	}

	newUser := db.User{
		Username: DEFAULT_ADMIN_USERNAME,
		Password: encoded,
	}

	err = s.db.CreateUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) initializeTokenSecret() (string, error) {
	existingSecret, err := s.db.GetConfig("jwt_secret")

	// secret already exists
	if err == nil {
		return existingSecret, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	bytes := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	secret := base64.URLEncoding.EncodeToString(bytes)

	err = s.db.InsertConfig("jwt_secret", secret)
	if err != nil {
		return "", err
	}

	return secret, nil
}
