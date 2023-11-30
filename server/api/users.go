package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) loginUserHandler(c echo.Context) error {
	fmt.Println("LOGIN USER")

	return nil
}
