// udpxy/udpxy.go
package udpxy

import (
	"net"
	"time"
)

type Udpxy struct {
	InterfaceName string
	Timeout       string
	inteface      *net.Interface
	timeout       time.Duration
}

// Provision prepares the udpxy instance for use.
func (u *Udpxy) Provision() error {
	inf, err := net.InterfaceByName(u.InterfaceName)
	if err != nil {
		return err
	}
	u.inteface = inf
	timeout, err := time.ParseDuration(u.Timeout)
	if err != nil {
		return err
	}
	u.timeout = timeout
	return nil
}
