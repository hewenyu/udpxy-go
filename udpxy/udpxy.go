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

// save interface
func (u *Udpxy) SaveInterface(i *net.Interface) {
	u.inteface = i
}

// save interface
func (u *Udpxy) SaveTimeout(t time.Duration) {
	u.timeout = t
}

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
