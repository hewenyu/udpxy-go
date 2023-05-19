package udp

import (
	"fmt"
	"testing"
)

func TestParseAddress(t *testing.T) {
	// parse address
	addr, err := parseAddress("igmp://233.50.201.133:5140")
	if err != nil {
		t.Fatal(err)
	}
	if addr == "" {
		t.Fatal("invalid address")
	}

	fmt.Println(addr)

}
