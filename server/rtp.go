package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// command like rtsp://124.75.34.37/PLTV/88888888/224/3221226078/00000100000000060000000000000321_0.smil
// use github.com/pion/rtp to parse rtp packets
func RTPHandler(c *gin.Context) {
	rtpURL := c.Param("command")

	fmt.Println(rtpURL)

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)

}
