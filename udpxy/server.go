package udpxy

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pion/rtp"
)

func (u *Udpxy) Serve(c *gin.Context) {
	parts := strings.Split(c.Request.URL.Path, "/")
	if len(parts) < 3 || parts[1] != "udp" {
		c.String(400, "No address specified")
		return
	}

	raddr := parts[2]
	addr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	conn, err := net.ListenMulticastUDP("udp4", u.inteface, addr)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add((u.timeout)))
	var buf = make([]byte, 1500)
	n, err := conn.Read(buf)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	conn.SetReadDeadline(time.Time{})
	p := &rtp.Packet{}
	headerSent := false
	for {
		if err = p.Unmarshal(buf[:n]); err != nil {
			c.String(500, err.Error())
			return
		}

		if !headerSent {
			headerSent = true
			if p.PayloadType == RTP_Payload_MP2T {
				c.Writer.Header().Set("Content-Type", ContentType_MP2T)
			} else {
				c.Writer.Header().Set("Content-Type", ContentType_DEFAULT)
			}
			c.Writer.WriteHeader(200)
		}

		if _, werr := c.Writer.Write(p.Payload); werr != nil {
			break
		}

		if n, err = conn.Read(buf); err != nil {
			break
		}
	}
	if err != nil && err != io.EOF {
		c.String(500, err.Error())
		return
	}
}
