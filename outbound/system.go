package outbound

import (
	"io"
	"net"
)

var SystemOutbound Outbound = &systemOutbound{}

type systemOutbound struct {
	net.Dialer
}

func (o *systemOutbound) DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	return net.DialTCP(network, laddr, raddr)
}

func (o *systemOutbound) Handshake(_ io.Reader, _ io.Writer, _, _ string) error {
	return nil
}
