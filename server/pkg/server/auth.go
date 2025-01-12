package server

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * JWT_EXPIRATION_TIME)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtConfig.SigningKey)
	if err != nil {
		s.logger.Error("could not create jwt", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

func (s *Server) refreshTokenHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		s.logger.Warn("authorization header is missing")
		return echo.NewHTTPError(http.StatusUnauthorized, "authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		s.logger.Warn("invalid authorization header format")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtConfig.SigningKey, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		s.logger.Errorln("could not parse token", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	claims, ok := parsedToken.Claims.(*jwtCustomClaims)
	if !ok {
		s.logger.Errorln("invalid claims type")
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	var refreshTokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	err = c.Bind(&refreshTokenRequest)
	if err != nil {
		s.logger.Warn("could not bind request json body", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	refreshTokenHash := sha512.Sum512([]byte(refreshTokenRequest.RefreshToken))
	refreshTokenHashString := hex.EncodeToString(refreshTokenHash[:])

	_, err = s.db.GetRefreshToken(claims.UserID, refreshTokenHashString)
	if err != nil {
		s.logger.Errorln("could not find refresh token", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	newRefreshToken := uuid.NewString()

	newRefreshTokenHash := sha512.Sum512([]byte(newRefreshToken))
	newRefreshTokenHashString := hex.EncodeToString(newRefreshTokenHash[:])

	err = s.db.UpdateRefreshToken(claims.UserID, refreshTokenHashString, newRefreshTokenHashString)
	if err != nil {
		s.logger.Errorln("could not update refresh token", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	newClaims := &jwtCustomClaims{
		claims.UserID,
		claims.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * JWT_EXPIRATION_TIME)),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	accessToken, err := newToken.SignedString(s.jwtConfig.SigningKey)
	if err != nil {
		s.logger.Error("could not create jwt", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": accessToken, "refresh_token": newRefreshToken})
}
