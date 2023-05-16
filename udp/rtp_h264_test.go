package udp

import (
	"log"
	"testing"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/formats/rtph264"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
)

func TestRTSP(t *testing.T) {

	c := gortsplib.Client{}

	u, err := url.Parse("rtsp://124.75.34.37/PLTV/88888888/224/3221226078/00000100000000060000000000000321_0.smil")
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

	// setup H264 -> raw frames decoder
	h264RawDec, err := newH264Decoder()
	if err != nil {
		t.Fatal(err)
	}
	defer h264RawDec.close()

	// if SPS and PPS are present into the SDP, send them to the decoder
	if forma.SPS != nil {
		h264RawDec.decode(forma.SPS)
	}
	if forma.PPS != nil {
		h264RawDec.decode(forma.PPS)
	}

	// setup a single media
	_, err = c.Setup(medi, baseURL, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// called when a RTP packet arrives
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// extract access units from RTP packets
		au, pts, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
			}
			return
		}

		for _, nalu := range au {
			// convert NALUs into RGBA frames
			img, err := h264RawDec.decode(nalu)
			if err != nil {
				t.Fatal(err)
			}

			// wait for a frame
			if img == nil {
				continue
			}

			log.Printf("decoded frame with size %v and pts %v", img.Bounds().Max, pts)
		}
	})

	// start playing
	_, err = c.Play(nil)
	if err != nil {
		t.Fatal(err)
	}

	// wait until a fatal error
	t.Fatal(c.Wait())
}
