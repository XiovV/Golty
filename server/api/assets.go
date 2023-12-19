package api

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) serveThumbnail(c echo.Context) error {
	thumbnail := c.Param("thumbnail")

	path := fmt.Sprintf("thumbnails/%s", thumbnail)
	return c.File(path)
}

func (s *Server) serveAvatar(c echo.Context) error {
	thumbnail := strings.Replace(c.Param("avatar"), "%40", "@", 1)

	path := fmt.Sprintf("avatars/%s", thumbnail)
	return c.File(path)
}
