package udpxy

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Register "udpxy" directive
func init() {
	httpcaddyfile.RegisterHandlerDirective("udpxy", parseCaddyfile)
}

// parseCaddyfile sets up the Udpxy middleware from Caddyfile
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	u := new(Udpxy)
	err := u.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//	udpxy {
//	    interface <interface>
//	    timeout <timeout>
//	}
func (u *Udpxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "interface":
				if d.NextArg() {
					u.InterfaceName = d.Val()
				}
			case "timeout":
				if d.NextArg() {
					u.Timeout = d.Val()
				}
			}
		}
	}
	return nil
}

// Assert Unmarshaler interface implementation
var (
	_ caddyfile.Unmarshaler = (*Udpxy)(nil)
)
