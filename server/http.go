// start http server
package server

import (
	"io"
	"log"
	"net/http"

	"github.com/hewenyu/udpxy-go/utils"
)

// HTTPServer 是一个结构，它包含一个HTTP服务器和一个数据channel
type HTTPServer struct {
	server      *http.Server // HTTP服务器
	dataChannel chan []byte  // 数据channel
}

// Start 方法启动HTTP服务器，该服务器从channel读取数据并将其写入HTTP响应
func (h *HTTPServer) Start(address string) error {
	// 创建一个channelReader
	reader := utils.NewChannelReader(h.dataChannel)

	// 创建一个HTTP服务器
	h.server = &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 设置响应头
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)

			// 将channel的数据复制到HTTP响应
			n, err := io.Copy(w, reader)
			if err != nil {
				// handle error
				log.Fatal(err)
			}
			log.Printf("%s %s %d [%s]", r.RemoteAddr, r.URL.Path, n, r.UserAgent())
		}),
	}

	// 在goroutine中启动HTTP服务器
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			// handle error
			log.Fatal(err)
		}
	}()

	return nil
}

// new
func NewHTTPServer(dataChannel chan []byte) *HTTPServer {
	return &HTTPServer{
		dataChannel: dataChannel,
	}
}
