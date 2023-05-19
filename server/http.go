package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hewenyu/udpxy-go/udp"
	"github.com/hewenyu/udpxy-go/utils"
)

// HTTPServer is a HTTP server that serves MPEG-TS over HTTP
type HTTPServer struct {
	server *http.Server
	pool   *sync.Map
	sem    chan struct{} // a semaphore for limiting the number of active connections
}

// Start listens on the specified address and starts serving HTTP requests
func (h *HTTPServer) Start(address string, maxConnections int) error {
	// create a semaphore with the specified size
	h.sem = make(chan struct{}, maxConnections)

	h.server = &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// try to acquire a slot from the semaphore
			select {
			case h.sem <- struct{}{}:
			default:
				// all slots are occupied, respond with 503 Service Unavailable
				http.Error(w, "Too many connections", http.StatusServiceUnavailable)
				return
			}

			// make sure to release the slot at the end
			defer func() { <-h.sem }()

			// r.URL.Path should be /udp/233.50.201.118:5140
			udpAddr := strings.TrimPrefix(r.URL.Path, "/udp/")
			connKey := r.RemoteAddr + udpAddr + time.Now().Format(time.RFC3339Nano) // use combination of RemoteAddr, udpAddr and timestamp as key

			log.Println("udpAddr:", udpAddr)

			ch := make(chan []byte, 100)
			h.pool.Store(connKey, ch) // store the channel with connKey in the pool

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			udpReceiver := udp.NewUDPReceiver(ch)
			if err := udpReceiver.Start(ctx, "eth0", udpAddr); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// make sure to stop the UDP receiver when the HTTP request is done
			go func() {
				<-r.Context().Done()
				udpReceiver.Stop()
				h.pool.Delete(connKey) // delete the channel with connKey from the pool
				close(ch)              // close the channel
			}()

			reader := utils.NewChannelReader(ch)

			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			n, err := io.Copy(w, reader)
			if err != nil {
				// handle error
				log.Println(err)
			}
			log.Printf("%s %s %d [%s]", r.RemoteAddr, r.URL.Path, n, r.UserAgent())
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

// NewHTTPServer creates a new HTTPServer instance
func NewHTTPServer(pool *sync.Map) *HTTPServer {
	return &HTTPServer{
		pool: pool,
	}
}
