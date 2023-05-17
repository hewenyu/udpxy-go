package utils

import "io"

// ChannelReader is a reader that reads from a channel
type ChannelReader struct {
	ch chan []byte // channel to read from
}

// Read channel reader
func (cr *ChannelReader) Read(p []byte) (n int, err error) {
	data := <-cr.ch
	n = copy(p, data)
	if n < len(data) {
		err = io.ErrShortBuffer
	}
	return n, err
}

// new
func NewChannelReader(ch chan []byte) *ChannelReader {
	return &ChannelReader{ch: ch}
}
