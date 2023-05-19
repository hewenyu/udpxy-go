package udp

import (
	"context"
	"log"
	"net"
	"sync"
)

// UDPReceiver is a UDP receiver
type UDPReceiver struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	conn       *net.UDPConn
	pool       *sync.Map
}

// Start listens for UDP packets on the specified interface and multicast address
func (u *UDPReceiver) Start(ctx context.Context, interfaceName string, multicastAddress string) error {
	u.ctx, u.cancelFunc = context.WithCancel(ctx)

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}
	// address, err := parseAddress(multicastAddress)
	// if err != nil {
	// 	return err
	// }

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
			select {
			case <-u.ctx.Done():
				return
			default:
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
		}
	}()

	return nil
}

// Stop stops the UDP receiver
func (u *UDPReceiver) Stop() error {
	if u.conn != nil {
		return u.conn.Close()
	}
	return nil
}

// NewUDPReceiver creates a new UDPReceiver instance
func NewUDPReceiver(pool *sync.Map) *UDPReceiver {
	return &UDPReceiver{
		pool: pool,
	}
}
