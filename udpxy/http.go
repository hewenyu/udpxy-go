package udpxy

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/pion/rtp"
)

// start http
func (u *Udpxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//
	parts := strings.FieldsFunc(r.URL.Path, func(r rune) bool { return r == '/' })
	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "No address specified")
		return next.ServeHTTP(w, r)
	}
	raddr := parts[1]
	addr, err := net.ResolveUDPAddr("udp4", raddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Invalid address specified")
		return next.ServeHTTP(w, r)
	}
	conn, err := net.ListenMulticastUDP("udp4", u.inteface, addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return next.ServeHTTP(w, r)
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(u.timeout))
	var buf = make([]byte, 1500)
	n, err := conn.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return next.ServeHTTP(w, r)
	}
	conn.SetReadDeadline(time.Time{})
	p := &rtp.Packet{}
	headerSent := false
	wc := int64(0) // Initialize write counter
	for {
		if err = p.Unmarshal(buf[:n]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return next.ServeHTTP(w, r)
		}

		if !headerSent {
			headerSent = true
			if p.PayloadType == RTP_Payload_MP2T {
				w.Header().Set("Content-Type", ContentType_MP2T)
			} else {
				w.Header().Set("Content-Type", ContentType_DEFAULT)
			}
			w.WriteHeader(http.StatusOK)
		}

		if _, werr := w.Write(p.Payload); werr != nil {
			break
		} else {
			wc += int64(n) // Update write counter
		}

		if n, err = conn.Read(buf); err != nil {
			break
		}
	}
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return next.ServeHTTP(w, r)
	}
	return next.ServeHTTP(w, r)
}

// Assert MiddlewareHandler interface implementation
var (
	_ caddyhttp.MiddlewareHandler = (*Udpxy)(nil)
)
