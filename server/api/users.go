package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) loginUserHandler(c echo.Context) error {
	userToken := s.getUserToken(c)
	fmt.Println(userToken)
	return nil
}
