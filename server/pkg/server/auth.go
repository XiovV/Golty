package server

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/matthewhartstonge/argon2"
)

type jwtCustomClaims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *Server) loginHandler(c echo.Context) error {
	var loginRequest struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	err := c.Bind(&loginRequest)
	if err != nil {
		s.logger.Warn("could not bind request json body", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	user, err := s.db.GetUserByUsername(loginRequest.Username)
	if err != nil {
		s.logger.Error("could not get user by username", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	ok, err := argon2.VerifyEncoded([]byte(loginRequest.Password), user.Password)
	if err != nil {
		s.logger.Error("could not verify password via argon2", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	if !ok {
		s.logger.Warn("invalid login credentials")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	refreshToken := uuid.NewString()

	refreshTokenHash := sha512.Sum512([]byte(refreshToken))
	refreshTokenHashString := hex.EncodeToString(refreshTokenHash[:])

	err = s.db.InsertRefreshToken(user.ID, refreshTokenHashString)
	if err != nil {
		s.logger.Error("could not insert refresh token", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	claims := &jwtCustomClaims{
		user.ID,
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		s.logger.Error("could not create jwt", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": accessToken, "refresh_token": refreshToken})
}
