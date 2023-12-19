package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) serveThumbnail(c echo.Context) error {
	thumbnail := c.Param("thumbnail")

	path := fmt.Sprintf("thumbnails/%s", thumbnail)
	return c.File(path)
}
