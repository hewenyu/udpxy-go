package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) handleStatus(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusOK, gin.H{
		"status": "Status information here...",
	})
}
