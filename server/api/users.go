package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) loginUserHandler(c echo.Context) error {
	var loginUserRequest struct {
		Username string `json:"username" validate:"required,lte=100"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginUserRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	validate := validator.New()

	err := validate.Struct(loginUserRequest)
	if err != nil {
		errors := strings.Split(err.Error(), "\n")

		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": errors})
	}

	user, err := s.Repository.FindUserByUsername(loginUserRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect username or password")
	}

	ok, err := argon2id.ComparePasswordAndHash(loginUserRequest.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect username or password")
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect username or password")
	}

	isAdmin := false
	if user.Admin == 1 {
		isAdmin = true
	}

	token, err := s.generateNewToken(user.ID, user.Username, isAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    fmt.Sprintf("Bearer %s", token),
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
	}

	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}
