package main

import (
	"fmt"

	"github.com/hewenyu/udpxy-go/server"
	"github.com/hewenyu/udpxy-go/udp"
)

func main() {
	// 创建一个数据channel
	dataChannel := make(chan []byte, 100)

	// 创建一个UDPReceiver
	udpReceiver := udp.NewUDPReceiver(dataChannel)

	// 创建一个HTTPServer
	httpServer := server.NewHTTPServer(dataChannel)

	// 启动UDPReceiver
	err := udpReceiver.Start("eth0", "224.0.0.1:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 启动HTTPServer
	err = httpServer.Start("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 阻止主函数退出
	select {}
}
