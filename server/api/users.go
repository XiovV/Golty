package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *Server) loginUserHandler(c *gin.Context) {
	fmt.Println("LOGIN USER")
}
