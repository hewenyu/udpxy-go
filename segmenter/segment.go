package segmenter

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/asticode/go-astits"
)

type Segment struct {
	Index    int
	Duration time.Duration
	FilePath string
}

type Segmenter struct {
	mux                   sync.Mutex
	segments              []Segment
	maxSegments           int           // The maximum number of segments to keep.
	currentSequenceNumber int           // current sequence number
	duration              time.Duration // segment duration
}

// Assuming there's a function to get the current sequence number.
func (s *Segmenter) getCurrentSequenceNumber() int {
	// Return the current sequence number...
	return s.currentSequenceNumber
}

// Assuming there's a function to get the segment duration.
func (s *Segmenter) getSegmentDuration() time.Duration {
	// Return the segment duration...
	return s.duration
}

// Assuming there's a function to get the segment file path.
func (s *Segmenter) getSegmentFilePath(sequenceNumber int) string {
	// Return the file path for the segment with the given sequence number...
	return s.segments[sequenceNumber].FilePath
}

func (s *Segmenter) UpdateSegments(conn *net.UDPConn) error {
	// Lock mutex while updating segments.
	s.mux.Lock()
	defer s.mux.Unlock()

	// Check if there are any expired segments.
	if len(s.segments) > s.maxSegments {
		// Remove the expired segment.
		expiredSegment := s.segments[0]
		s.segments = s.segments[1:]
		// Delete the expired segment file.
		err := os.Remove(expiredSegment.FilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open the .ts file
	f, err := os.Open("/path/to/your/file.ts")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a cancellable context in case you want to stop reading packets/data any time you want
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all paths cancel the context to prevent context leak

	// Create the demuxer with the UDP connection as the input.
	dmx := astits.NewDemuxer(ctx, conn)

	// Create a new segment.
	segment := Segment{
		Index:    s.getCurrentSequenceNumber(),
		Duration: s.getSegmentDuration(),
		FilePath: s.getSegmentFilePath(s.getCurrentSequenceNumber()),
	}
	s.segments = append(s.segments, segment)

	// Create a new file for the segment.
	out, err := os.Create(segment.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Loop through data
	var pd *astits.DemuxerData
	var perr error

	for {
		// Get the next data
		pd, perr = dmx.NextData()
		if perr != nil {
			// Break the loop if an error occurs
			log.Println(perr)
			break
		}

		// Data is a PMT data
		if pd.PMT != nil {
			// Loop through elementary streams
			for _, es := range pd.PMT.ElementaryStreams {
				fmt.Printf("Stream detected: %d\n", es.ElementaryPID)
			}
		}

		_, perr = out.Write(pd.FirstPacket.Payload)
		if err != nil {
			log.Fatal(perr)
		}

		// ...existing code to write packet data to the output file...

		// TODO: Handle segmenting based on timecodes, update M3U8 playlist, manage segment indices and expirations...
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return nil
}
