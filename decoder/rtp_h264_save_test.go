package decoder

import (
	"log"
	"testing"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/formats/rtph264"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
)

// get real-time stream
// search for the public rtsp stream

const TestUrl = "rtsp://rtspstream.com/pattern"

func TestSaveDisk(t *testing.T) {
	c := gortsplib.Client{}

	u, err := url.Parse(TestUrl)
	if err != nil {
		t.Fatal(err)
	}

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// find published medias
	medias, baseURL, _, err := c.Describe(u)
	if err != nil {
		t.Fatal(err)
	}

	// find the H264 media and format
	var forma *formats.H264
	medi := medias.FindFormat(&forma)
	if medi == nil {
		t.Fatal("media not found")
	}

	// setup RTP/H264 -> H264 decoder
	rtpDec := forma.CreateDecoder()

	// setup H264 -> MPEG-TS muxer
	mpegtsMuxer, err := newMPEGTSMuxer(forma.SPS, forma.PPS, "test.ts")
	if err != nil {
		t.Fatal(err)
	}

	// setup a single media
	_, err = c.Setup(medi, baseURL, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// called when a RTP packet arrives
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// extract access unit from RTP packets
		// DecodeUntilMarker is necessary for the DTS extractor to work
		au, pts, err := rtpDec.DecodeUntilMarker(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
			}
			return
		}

		// encode the access unit into MPEG-TS
		mpegtsMuxer.encode(au, pts)
	})

	// start playing
	_, err = c.Play(nil)
	if err != nil {
		t.Fatal(err)
	}

	// wait until a fatal error
	t.Fatal(c.Wait())
}
