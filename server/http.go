// start http server
package server

import (
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router *gin.Engine
}

func NewHTTPServer() *HTTPServer {
	router := gin.Default()

	httpServer := &HTTPServer{
		router: router,
	}

	router.GET("/rtp/:command", RTPHandler)
	router.GET("/status", httpServer.handleStatus)

	return httpServer
}

func (s *HTTPServer) Start(addr string) {
	s.router.Run(addr) // listen and serve on
}
