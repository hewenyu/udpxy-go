package udp

import (
	"log"
	"net"
	"sync"
)

// UDPReceiver is a UDP receiver
type UDPReceiver struct {
	conn *net.UDPConn
	pool *sync.Map
}

// Start listens for UDP packets on the specified interface and multicast address
func (u *UDPReceiver) Start(interfaceName string, multicastAddress string) error {
	// interfaceName is a string like "eth0"
	// multicastAddress is a string like "igmp://233.50.201.133:5140"

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	// 使用自定义解析函数
	address, err := parseAddress(multicastAddress)
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return err
	}

	u.conn, err = net.ListenMulticastUDP("udp4", iface, addr)
	if err != nil {
		return err
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := u.conn.ReadFromUDP(buffer)
			if err != nil {
				// handle error
				log.Println(err)
				continue
			}

			u.pool.Range(func(key, value interface{}) bool {
				ch := value.(chan []byte)
				select {
				case ch <- buffer[:n]:
				default:
					// channel is full, drop the packet
				}
				return true
			})
		}
	}()

	return nil
}

// new
func NewUDPReceiver(pool *sync.Map) *UDPReceiver {
	return &UDPReceiver{
		pool: pool,
	}
}
