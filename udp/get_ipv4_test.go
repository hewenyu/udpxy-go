package udp

// test for GetIPv4
import (
	"fmt"
	"testing"
)

func TestGetIPv4(t *testing.T) {
	// get IPv4 address
	ip, err := getIPv4("eth0")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("IPv4 address:", ip)
}
