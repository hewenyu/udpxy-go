package udp

import (
	"github.com/bluenviron/gortsplib/v3/pkg/url"
)

func parseAddress(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	// 返回没有协议头的地址
	return u.Host, nil
}

// func parseAddress(raw string) (string, error) {
// 	if !strings.HasPrefix(raw, "igmp://") {
// 		return "", fmt.Errorf("invalid protocol")
// 	}

// 	// remove the protocol
// 	addr := strings.TrimPrefix(raw, "igmp://")

// 	// check that the address is not empty after removing the protocol
// 	if len(addr) == 0 {
// 		return "", fmt.Errorf("invalid address")
// 	}

// 	return addr, nil
// }
