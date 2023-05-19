package udpxy

import (
	"fmt"
	"net"
	"time"

	"github.com/caddyserver/caddy/v2"
)

// Register Udpxy module on init
func init() {
	caddy.RegisterModule(Udpxy{})
}

// Udpxy struct with configuration fields and internal fields
type Udpxy struct {
	InterfaceName string `json:"interface"`
	Timeout       string `json:"timeout"`
	inteface      *net.Interface
	timeout       time.Duration
}

// CaddyModule returns the Caddy module information.
func (Udpxy) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.udpxy",
		New: func() caddy.Module { return new(Udpxy) },
	}
}

// Validate checks if interface is not nil
func (u *Udpxy) Validate() error {
	if u.inteface == nil {
		return fmt.Errorf("no interface")
	}
	return nil
}

// Provision initializes and validates the necessary fields in the Udpxy structure
func (u *Udpxy) Provision(ctx caddy.Context) error {
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

// Assert interface implementations
var (
	_ caddy.Provisioner = (*Udpxy)(nil)
	_ caddy.Validator   = (*Udpxy)(nil)
)
