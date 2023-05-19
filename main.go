package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var addr = flag.String("l", ":18000", "Listening address")
var iface = flag.String("i", "eth0", "Listening multicast interface")
var inf *net.Interface

type UDPReceiver struct {
	conn *net.UDPConn
}

func NewUDPReceiver(address string) (*UDPReceiver, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp4", inf, addr)
	if err != nil {
		return nil, err
	}

	return &UDPReceiver{
		conn: conn,
	}, nil
}

func (u *UDPReceiver) Close() error {
	return u.conn.Close()
}

func (u *UDPReceiver) CopyTo(w io.Writer) (int64, error) {
	return io.Copy(w, u.conn)
}

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer() *HTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/rtp/", handleHTTP)
	return &HTTPServer{
		mux: mux,
	}
}

func (h *HTTPServer) Start(address string) error {
	return http.ListenAndServe(address, h.mux)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	parts := strings.FieldsFunc(req.URL.Path, func(r rune) bool {
		return r == '/'
	})

	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "No address specified")
		return
	}

	raddr := parts[1]
	receiver, err := NewUDPReceiver(raddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	defer receiver.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	n, err := receiver.CopyTo(w)
	if err != nil {
		log.Printf("ERR: %v", err)
		return
	}
	log.Printf("%s %s %d [%s]", req.RemoteAddr, req.URL.Path, n, req.UserAgent())
}

func main() {
	if os.Getppid() == 1 {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	} else {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	}

	flag.Parse()

	var err error
	inf, err = net.InterfaceByName(*iface)
	if err != nil {
		log.Fatal(err)
		return
	}

	server := NewHTTPServer()

	log.Fatal(server.Start(*addr))
}
