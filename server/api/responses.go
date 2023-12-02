package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Error string
}

func (s *Server) badRequestResponse(c echo.Context, message string) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: message})
}
