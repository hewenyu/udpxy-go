// for receiver udp packets
package udp

type Receiver interface {
	Start(interfaceName string, multicastAddress string) error // Start listening for incoming packets
}
