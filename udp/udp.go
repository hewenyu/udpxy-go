package udp

import (
	"log"
	"net"
)

// UDPReceiver 是一个结构，它包含一个UDP连接和一个数据channel
type UDPReceiver struct {
	conn        *net.UDPConn // UDP连接
	dataChannel chan []byte  // 数据channel
}

// Start 方法在指定的网络接口和多播地址上启动UDP接收器
func (u *UDPReceiver) Start(interfaceName string, multicastAddress string) error {
	// 通过名称获取网络接口
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	// 解析多播地址
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil {
		return err
	}

	// 在网络接口上监听多播UDP
	u.conn, err = net.ListenMulticastUDP("udp4", iface, addr)
	if err != nil {
		return err
	}

	// 在goroutine中读取UDP数据并将其写入channel
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := u.conn.ReadFromUDP(buffer)
			if err != nil {
				// handle error
				log.Fatal(err)
			}
			u.dataChannel <- buffer[:n]
		}
	}()

	return nil
}

// new
func NewUDPReceiver(dataChannel chan []byte) *UDPReceiver {
	return &UDPReceiver{
		dataChannel: dataChannel,
	}
}
