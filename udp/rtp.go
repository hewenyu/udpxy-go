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

type RTPReceiver struct {
	conn *gortsplib.Client
	pool *sync.Map
}

// Start listening for incoming packets
func (r *RTPReceiver) Start(interfaceName string, rtspUrl string) error {
	// find out the local IP of the interface
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	// use the first IPv4 address
	var localIP net.IP
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		if ip.To4() != nil {
			localIP = ip
			break
		}
	}

	if localIP == nil {
		return fmt.Errorf("no IPv4 address found for interface %s", interfaceName)
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

	err = client.Start(u.Scheme, u.Host)
	if err != nil {
		return err
	}

	medias, _, _, err := client.Describe(u)
	if err != nil {
		return err
	}

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

func NewRTPReceiver(pool *sync.Map) *RTPReceiver {
	return &RTPReceiver{
		pool: pool,
	}
}
