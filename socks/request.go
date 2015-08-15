package socks

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	connectCommand   = uint8(1)
	bindCommand      = uint8(2)
	associateCommand = uint8(3)
	ipv4Address      = uint8(1)
	dnsAddress       = uint8(3)
	ipv6Address      = uint8(4)
)

func (p *Proxy) HandleRequest(conn io.Writer, bufConn io.Reader) error {
	header := []byte{0, 0, 0}
	if _, err := io.ReadAtLeast(bufConn, header, 3); err != nil {
		return fmt.Errorf("Failed to get command version: %v", err)
	}

	if header[0] != socks5Version {
		return fmt.Errorf("Unsupported command version: %v", header[0])
	}

	dest, err := p.ReadAddress(conn, bufConn)
	if err != nil {
		p.Log.Errorf("Failed to read address: %v", err)
		return err
	}

	switch header[1] {
	case connectCommand:
		return p.HandleConnect(conn, bufConn, dest)
	case bindCommand:
		return nil
	case associateCommand:
		return nil
	default:
		return fmt.Errorf("Unsupported command: %v", header[1])
	}
}

type Address struct {
	IP   net.IP
	Port int
}

func (p *Proxy) ReadAddress(conn io.Writer, bufConn io.Reader) (*Address, error) {
	addrType := []byte{0}
	if _, err := bufConn.Read(addrType); err != nil {
		return nil, err
	}

	address := &Address{}

	switch addrType[0] {
	case ipv4Address:
		addr := make([]byte, 4)
		if _, err := io.ReadAtLeast(bufConn, addr, len(addr)); err != nil {
			return nil, err
		}
		address.IP = net.IP(addr)

	case ipv6Address:
		addr := make([]byte, 16)
		if _, err := io.ReadAtLeast(bufConn, addr, len(addr)); err != nil {
			return nil, err
		}
		address.IP = net.IP(addr)

	case dnsAddress:
		if _, err := bufConn.Read(addrType); err != nil {
			return nil, err
		}
		addrLen := int(addrType[0])
		domainName := make([]byte, addrLen)
		if _, err := io.ReadAtLeast(bufConn, domainName, addrLen); err != nil {
			return nil, err
		}
		addr, err := net.ResolveIPAddr("ip", string(domainName))
		if err != nil {
			return nil, err
		}
		address.IP = addr.IP

	default:
		// TODO respond back
		return nil, fmt.Errorf("Unknown address type: %q", addrType[0])
	}

	port := []byte{0, 0}
	if _, err := io.ReadAtLeast(bufConn, port, 2); err != nil {
		return nil, err
	}
	address.Port = (int(port[0]) << 8) | int(port[1])

	return address, nil
}

func (p *Proxy) HandleConnect(
	conn io.Writer,
	bufConn io.Reader,
	dest *Address) error {

	addr := net.TCPAddr{IP: dest.IP, Port: dest.Port}
	target, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		return err
	}
	defer target.Close()

	// Send success
	local := target.LocalAddr().(*net.TCPAddr)
	bind := &Address{IP: local.IP, Port: local.Port}
	if err := p.SendReply(conn, successReply, bind); err != nil {
		return fmt.Errorf("Failed to send reply: %v", err)
	}

	// Start proxying
	errCh := make(chan error, 2)
	go proxy("target", target, bufConn, errCh)
	go proxy("client", conn, target, errCh)

	// Wait
	select {
	case e := <-errCh:
		return e
	}
}

func (p *Proxy) SendReply(w io.Writer, resp uint8, addr *Address) error {
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case addr == nil:
		addrType = ipv4Address
		addrBody = []byte{0, 0, 0, 0}
		addrPort = 0

	case addr.IP.To4() != nil:
		addrType = ipv4Address
		addrBody = []byte(addr.IP.To4())
		addrPort = uint16(addr.Port)

	case addr.IP.To16() != nil:
		addrType = ipv6Address
		addrBody = []byte(addr.IP.To16())
		addrPort = uint16(addr.Port)

	default:
		return fmt.Errorf("Failed to format address: %v", addr)
	}

	msg := make([]byte, 6+len(addrBody))
	msg[0] = socks5Version
	msg[1] = resp
	msg[2] = 0
	msg[3] = addrType
	copy(msg[4:], addrBody)
	msg[4+len(addrBody)] = byte(addrPort >> 8)
	msg[4+len(addrBody)+1] = byte(addrPort & 0xff)

	_, err := w.Write(msg)
	return err
}

// proxy is used to suffle data from src to destination, and sends errors
// down a dedicated channel
func proxy(name string, dst io.Writer, src io.Reader, errCh chan error) {
	// Copy
	n, err := io.Copy(dst, src)

	// Log, and sleep. This is jank but allows the otherside
	// to finish a pending copy
	log.Printf("[DEBUG] socks: Copied %d bytes to %s", n, name)
	time.Sleep(10 * time.Millisecond)

	// Send any errors
	errCh <- err
}
