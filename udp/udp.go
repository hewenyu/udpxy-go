package udp

import (
	"context"
	"log"
	"net"
)

// UDPReceiver is a UDP receiver
type UDPReceiver struct {
	ctx    context.Context
	conn   *net.UDPConn
	output chan []byte
}

// Start listens for UDP packets on the specified interface and multicast address
func (u *UDPReceiver) Start(ctx context.Context, interfaceName string, multicastAddress string) error {
	u.ctx = ctx

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}
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
				log.Println(err)
				continue
			}

			select {
			case <-u.ctx.Done():
				return
			case u.output <- buffer[:n]: // send the packet to the output channel
			default:
				// output channel is full, drop the packet
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
func NewUDPReceiver(output chan []byte) *UDPReceiver {
	return &UDPReceiver{
		output: output,
	}
}
