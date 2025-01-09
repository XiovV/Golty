package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthewhartstonge/argon2"
)

func (s *Server) loginHandler(c echo.Context) error {
	var loginRequest struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	err := c.Bind(&loginRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	user, err := s.db.GetUserByUsername(loginRequest.Username)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	ok, err := argon2.VerifyEncoded([]byte(loginRequest.Password), user.Password)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	if !ok {
		fmt.Println("invalid login credentials")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login credentials")
	}

	fmt.Println("login user")

	return nil
}
