package udp

import (
	"log"
	"net"
	"sync"
)

// UDPReceiver 是一个结构，它包含一个UDP连接和一个数据channel
type UDPReceiver struct {
	conn *net.UDPConn
	pool *sync.Map
}

// Start listens for UDP packets on the specified interface and multicast address
func (u *UDPReceiver) Start(interfaceName string, multicastAddress string) error {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
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
