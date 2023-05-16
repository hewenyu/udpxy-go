package utils

import "io"

// ChannelReader 是一个结构，它包装了一个channel并实现了io.Reader接口
type ChannelReader struct {
	ch chan []byte // 数据channel
}

// Read 方法从channel读取数据并将其复制到提供的缓冲区
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
