package udp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
)

// RTPReceiver is a RTP receiver
type RTPReceiver struct {
	conn *gortsplib.Client
	pool *sync.Map
}

// Start listening for incoming packets
func (r *RTPReceiver) Start(interfaceName string, rtspUrl string) error {
	// find out the local IP of the interface
	localIP, err := getIPv4(interfaceName)
	if err != nil {
		return err
	}

	client := gortsplib.Client{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				LocalAddr: &net.TCPAddr{
					IP: localIP,
				},
			}

			return dialer.DialContext(ctx, network, address)
		},
	}

	u, err := url.Parse(rtspUrl)
	if err != nil {
		return err
	}
	// connect to the server
	err = client.Start(u.Scheme, u.Host)
	if err != nil {
		return err
	}

	medias, _, _, err := client.Describe(u)
	if err != nil {
		return err
	}

	// find the H264 media and format
	var forma *formats.H264
	medi := medias.FindFormat(&forma)
	if medi == nil {
		return fmt.Errorf("media not found")
	}

	client.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		r.pool.Range(func(key, value interface{}) bool {
			ch := value.(chan []byte)
			select {
			case ch <- pkt.Payload:
			default:
				// channel is full, drop the packet
			}
			return true
		})
	})

	_, err = client.Play(nil)
	if err != nil {
		return err
	}

	r.conn = &client
	return nil
}

// NewRTPReceiver creates a new RTPReceiver
func NewRTPReceiver(pool *sync.Map) *RTPReceiver {
	return &RTPReceiver{
		pool: pool,
	}
}
