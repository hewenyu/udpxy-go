package main

import (
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hewenyu/udpxy-go/udpxy"
)

func main() {
	router := gin.Default()

	u := &udpxy.Udpxy{
		InterfaceName: "eth0", // your network interface here
		Timeout:       "30s",  // your timeout here
	}

	inf, err := net.InterfaceByName(u.InterfaceName)
	if err != nil {
		log.Fatalf("error setting interface: %v", err)
	}
	u.SaveInterface(inf)
	timeout, err := time.ParseDuration(u.Timeout)
	if err != nil {
		log.Fatalf("error parsing duration: %v", err)
	}
	u.SaveTimeout(timeout)

	router.GET("/udp/:addr", u.Serve)

	router.Run(":9096")
}
