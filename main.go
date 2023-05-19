package main

import (
	"log"
	"sync"

	"github.com/hewenyu/udpxy-go/server"
)

func main() {
	// save udp address and channel
	pool := &sync.Map{}

	// max connections
	httpServer := server.NewHTTPServer(pool)
	err := httpServer.Start(":9096", 100)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("http server started")
	// block forever
	select {}
}
