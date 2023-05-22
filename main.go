package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hewenyu/udpxy-go/udpxy"
)

func main() {
	router := gin.Default()

	u := &udpxy.Udpxy{
		InterfaceName: "eth0", // your network interface here
		Timeout:       "30s",  // your timeout here
	}

	err := u.Provision()
	if err != nil {
		panic(err)
	}

	router.GET("/udp/:addr", u.Serve)

	// router.GET("/hls", segmenter.ServeHLS)
	// router.GET("/hls/:segment", segmenter.ServeSegment)

	router.Run(":9096")
}
