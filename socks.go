package wiretap

import (
	"net"

	"github.com/armon/go-socks5"
)

type SocksProxy struct {
	HttpPort  int
	HttpsPort int
	Server    *socks5.Server
}

func NewSocksProxy() (*SocksProxy, error) {
	s := SocksProxy{
		HttpPort:  9080,
		HttpsPort: 9443,
	}
	c := &socks5.Config{
		Rewriter: s,
	}
	var err error
	s.Server, err = socks5.New(c)
	return &s, err
}

func (s SocksProxy) ListenAndServe(network, addr string) error {
	return s.Server.ListenAndServe(network, addr)
}

// Rewrite makes SocksProxy implement the socks5.Rewriter interface
// Moves connection to local http proxy server
func (s SocksProxy) Rewrite(addr *socks5.AddrSpec) *socks5.AddrSpec {
	addr.IP = net.IP{0, 0, 0, 0}
	// TODO http or https
	addr.Port = s.HttpPort
	return addr
}
