// start http server
package server

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/hewenyu/udpxy-go/utils"
)

// HTTPServer 是一个结构，它包含一个HTTP服务器和一个数据channel
type HTTPServer struct {
	server *http.Server
	pool   *sync.Map
	sem    chan struct{} // a semaphore for limiting the number of active connections
}

// Start 方法启动HTTP服务器，该服务器从channel读取数据并将其写入HTTP响应
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

			ch := make(chan []byte, 100)
			h.pool.Store(r, ch)
			defer h.pool.Delete(r)

			reader := utils.NewChannelReader(ch)

			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			n, err := io.Copy(w, reader)
			if err != nil {
				// handle error
			}
			log.Printf("%s %s %d [%s]", r.RemoteAddr, r.URL.Path, n, r.UserAgent())
		}),
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			// handle error
		}
	}()

	return nil
}

// new
func NewHTTPServer(pool *sync.Map) *HTTPServer {
	return &HTTPServer{
		pool: pool,
	}
}
