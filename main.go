package main

import (
	"fmt"
	"sync"

	"github.com/hewenyu/udpxy-go/server"
	"github.com/hewenyu/udpxy-go/udp"
)

func main() {
	pool := &sync.Map{}

	// 创建一个UDPReceiver
	udpReceiver := udp.NewUDPReceiver(pool)

	// 创建一个HTTPServer
	httpServer := server.NewHTTPServer(pool)

	// 启动UDPReceiver
	err := udpReceiver.Start("eth0", "224.0.0.1:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 启动HTTPServer
	err = httpServer.Start("localhost:8080", 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 阻止主函数退出
	select {}
}
