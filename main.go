package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type UDPReceiver struct {
	conn        *net.UDPConn
	dataChannel chan []byte
}

func (u *UDPReceiver) Start(interfaceName string, port int) error {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	ip, _, err := net.ParseCIDR(addrs[0].String())
	if err != nil {
		return err
	}

	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}
	u.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := u.conn.ReadFromUDP(buffer)
			if err != nil {
				// handle error
				log.Fatal(err)
			}
			// Here you should handle RTP and MPEG-TS payloads
			u.dataChannel <- buffer[:n]
		}
	}()
	return nil
}

type HTTPServer struct {
	server      *http.Server
	dataChannel chan []byte
}

type channelReader struct {
	ch chan []byte
}

func (cr *channelReader) Read(p []byte) (n int, err error) {
	data := <-cr.ch
	n = copy(p, data)
	if n < len(data) {
		err = io.ErrShortBuffer
	}
	return n, err
}

func (h *HTTPServer) Start(address string) error {
	reader := &channelReader{ch: h.dataChannel}

	h.server = &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Here you should parse HTTP commands and act accordingly
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, reader)
		}),
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			// handle error
			log.Fatal(err)
		}
	}()
	return nil
}

// main function as before...
func main() {
	dataChannel := make(chan []byte, 100)

	udpReceiver := &UDPReceiver{
		dataChannel: dataChannel,
	}
	httpServer := &HTTPServer{
		dataChannel: dataChannel,
	}

	err := udpReceiver.Start("eth0", 12345)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = httpServer.Start("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prevent the main function from exiting
	select {}
}
