package main

import (
	"fmt"
	"log"
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
	err := udpReceiver.Start("eth0", "igmp://233.50.201.133:5140")
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("start udp success")

	// 启动HTTPServer
	err = httpServer.Start("localhost:9096", 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("start http success")

	// 阻止主函数退出
	select {}
}
