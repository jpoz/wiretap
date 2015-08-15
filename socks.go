package wiretap

import (
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/armon/go-socks5"
	"github.com/jpoz/wiretap/disk"
)

type SocksProxy struct {
	HttpPort  int
	HttpsPort int
	Server    *socks5.Server
	HttpProxy HttpProxy
}

func NewSocksProxy() (*SocksProxy, error) {
	s := SocksProxy{
		HttpPort:  9080,
		HttpsPort: 9443,
	}
	c := &socks5.Config{
		Rewriter: s,
	}

	tr := &Transport{
		Storage: disk.Storage{"./cache"},
	}

	// HTTP client with wiretap Transport
	client := &http.Client{Transport: tr}

	// Proxy
	s.HttpProxy = HttpProxy{
		Client:   client,
		Director: director,
	}

	var err error
	s.Server, err = socks5.New(c)
	return &s, err
}

// TODO config this
func (s SocksProxy) ListenAndServe() error {
	log.Infof("Starting API server on %s", ":8000")

	errChan := make(chan error)

	go func() {
		errChan <- http.ListenAndServe(":9080", s.HttpProxy.ServeHTTP)
	}()
	go func() {
		errChan <- s.Server.ListenAndServe("tcp", ":8000")
	}()

	return <-errChan
}

// Rewrite makes SocksProxy implement the socks5.Rewriter interface
// Moves connection to local http proxy server
func (s SocksProxy) Rewrite(addr *socks5.AddrSpec) *socks5.AddrSpec {
	log.Infof("%+v", addr)
	addr.IP = net.IP{0, 0, 0, 0}
	// TODO http or https
	addr.Port = s.HttpPort
	return addr
}
