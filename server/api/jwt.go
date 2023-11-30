package api

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtConfig = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(jwtCustomClaims)
	},
	SigningKey: []byte("secret"),
}

type UserToken struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

type jwtCustomClaims struct {
	UserToken
	jwt.RegisteredClaims
}

func (s *Server) getUserToken(c echo.Context) UserToken {
	authToken := c.Get("user").(*jwt.Token)
	claims := authToken.Claims.(*jwtCustomClaims)

	return UserToken{
		Id:       claims.Id,
		Username: claims.Username,
		Admin:    claims.Admin,
	}
}
