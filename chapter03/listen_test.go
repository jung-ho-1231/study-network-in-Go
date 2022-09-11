package chapter03

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		panic(err)
	}

	defer func() { listener.Close() }()

	t.Logf("bound to %q", listener.Addr())
}
