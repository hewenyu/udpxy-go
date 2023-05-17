package udp

import (
	"errors"
	"net"
)

// get ipv4 address by interface name
// Path: udp/get_ipv4.go
func getIPv4(ifaceName string) (net.IP, error) {
	// find out the local IP of the interface
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	// use the first IPv4 address
	var localIP net.IP
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		if ip.To4() != nil {
			localIP = ip
			break
		}
	}

	if localIP == nil {
		return nil, errors.New("no IPv4 address found for interface " + ifaceName)
	}

	return localIP, nil
}
