package socks

import (
	"bufio"
	"io"
	"net"

	"github.com/Sirupsen/logrus"
)

const (
	socks5Version = uint8(5)

	noAuthType  = uint8(0)
	authSuccess = uint8(0)
	authFailure = uint8(1)
)

const (
	successReply uint8 = iota
	serverFailure
	ruleFailure
	networkUnreachable
	hostUnreachable
	connectionRefused
	ttlExpired
	commandNotSupported
	addrTypeNotSupported
)

type Proxy struct {
	Network string
	Addr    string
	Log     *logrus.Logger
}

func NewProxy(network, addr string) (*Proxy, error) {
	proxy := &Proxy{
		Network: network,
		Addr:    addr,
		Log:     logrus.New(),
	}
	return proxy, nil
}

func (p *Proxy) ListenAndServe() error {
	listener, err := net.Listen(p.Network, p.Addr)
	if err != nil {
		return err
	}
	return p.Serve(listener)
}

func (p *Proxy) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go p.HandleConn(conn)
	}
	return nil
}

func (p *Proxy) HandleConn(conn net.Conn) {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		p.Log.Errorf("Could not get version byte: %v", err)
	}

	if version[0] != socks5Version {
		p.Log.Errorf("Unsupported SOCKS version: %v", version)
	}

	if err := p.Authenticate(conn, bufConn); err != nil {
		p.Log.Errorf("Failed to Authenticate: %v", err)
	}

	if err := p.HandleRequest(conn, bufConn); err != nil {
		p.Log.Errorf("Failed to handle request: %v", err)
	}
}

func (p *Proxy) Authenticate(conn io.Writer, bufConn io.Reader) error {
	header := []byte{0}
	if _, err := bufConn.Read(header); err != nil {
		return err
	}

	numMethods := int(header[0])
	methods := make([]byte, numMethods)
	if _, err := io.ReadAtLeast(bufConn, methods, numMethods); err != nil {
		return err
	}

	// Ignore Authentication types, just try NoAuth
	return p.NoAuthAuthenticate(conn, bufConn)
}

func (p *Proxy) NoAuthAuthenticate(conn io.Writer, bufConn io.Reader) error {
	_, err := conn.Write([]byte{socks5Version, noAuthType})
	return err
}
