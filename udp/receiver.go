package udp

type Receiver interface {
	Start() error // Start listening for incoming packets
}
