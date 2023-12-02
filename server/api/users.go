package api

import (
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/labstack/echo/v4"
)

func (s *Server) loginUserHandler(c echo.Context) error {
	var loginUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginUserRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
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

	return c.NoContent(http.StatusOK)
}
