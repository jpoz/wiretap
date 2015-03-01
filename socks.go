package wiretap

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	socks5Version = uint8(5)
)

const (
	connectCommand   = uint8(1)
	bindCommand      = uint8(2)
	associateCommand = uint8(3)
	ipv4Address      = uint8(1)
	fqdnAddress      = uint8(3)
	ipv6Address      = uint8(4)
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

const (
	noAuth          = uint8(0)
	noAcceptable    = uint8(255)
	userPassAuth    = uint8(2)
	userAuthVersion = uint8(1)
	authSuccess     = uint8(0)
	authFailure     = uint8(1)
)

var (
	unrecognizedAddrType = fmt.Errorf("Unrecognized address type")
)

// AddrSpec is used to return the target AddrSpec
// which may be specified as IPv4, IPv6, or a FQDN
type AddrSpec struct {
	FQDN string
	IP   net.IP
	Port int
}

type Socks struct {
}

func (s Socks) Start(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.ServeConnection(conn)
	}
	return nil
}

func (s Socks) ServeConnection(conn net.Conn) error {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	// Read the version byte
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		log.Printf("[ERR] socks: Failed to get version byte: %v", err)
		return err
	}

	fmt.Printf("version: %+v\n", version)

	// Ensure we are compatible
	if version[0] != socks5Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v", version)
		log.Printf("[ERR] socks: %v", err)
		return err
	}

	// Read in auth octets
	authHeader := []byte{0}
	if _, err := bufConn.Read(authHeader); err != nil {
		return err
	}
	numMethods := int(authHeader[0])
	methods := make([]byte, numMethods)
	if _, err := io.ReadAtLeast(bufConn, methods, numMethods); err != nil {
		return err
	}
	fmt.Printf("auth: %+v %+v\n", authHeader, methods)
	// TODO If client doesn't support No auth then err

	// Write back No Auth method respones
	bytes := []byte{socks5Version, noAuth}
	fmt.Println("<-", bytes)
	if _, err := conn.Write(bytes); err != nil {
		return err
	}

	// Process the client request
	if err := s.handleRequest(conn, bufConn); err != nil {
		err = fmt.Errorf("Failed to handle request: %v", err)
		log.Printf("[ERR] socks: %v", err)
		return err
	}

	return nil
}

func (s *Socks) handleRequest(conn net.Conn, bufConn io.Reader) error {
	// Read the version byte
	header := []byte{0, 0, 0}
	if _, err := io.ReadAtLeast(bufConn, header, 3); err != nil {
		return fmt.Errorf("Failed to get command version: %v", err)
	}

	fmt.Printf("header %+v\n", header)

	// Ensure we are compatible
	if header[0] != socks5Version {
		return fmt.Errorf("Unsupported command version: %v", header[0])
	}

	// Read in the destination address
	dest, err := readAddrSpec(bufConn)
	if err != nil {
		if err == unrecognizedAddrType {
			if err := sendReply(conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("Failed to send reply: %v", err)
			}
		}
		return fmt.Errorf("Failed to read destination address: %v", err)
	}

	// Resolve the address if we have a FQDN
	if dest.FQDN != "" {
		addr, err := net.ResolveIPAddr("ip", dest.FQDN)
		if err != nil {
			return err
		}

		if err != nil {
			if err := sendReply(conn, hostUnreachable, nil); err != nil {
				return fmt.Errorf("Failed to send reply: %v", err)
			}
			return fmt.Errorf("Failed to resolve destination '%v': %v", dest.FQDN, err)
		}
		dest.IP = addr.IP
	}

	// Switch on the command
	switch header[1] {
	case connectCommand:
		return s.handleConnect(conn, bufConn, dest, dest)
	default:
		if err := sendReply(conn, commandNotSupported, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v", err)
		}
		return fmt.Errorf("Unsupported command: %v", header[1])
	}
}

// handleConnect is used to handle a connect command
func (s *Socks) handleConnect(conn net.Conn, bufConn io.Reader, dest, realDest *AddrSpec) error {
	// Attempt to connect
	addr := net.TCPAddr{IP: realDest.IP, Port: realDest.Port}
	target, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		msg := err.Error()
		resp := hostUnreachable
		if strings.Contains(msg, "refused") {
			resp = connectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = networkUnreachable
		}
		if err := sendReply(conn, resp, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v", err)
		}
		return fmt.Errorf("Connect to %v failed: %v", dest, err)
	}
	defer target.Close()

	// Send success
	local := target.LocalAddr().(*net.TCPAddr)
	bind := AddrSpec{IP: local.IP, Port: local.Port}
	if err := sendReply(conn, successReply, &bind); err != nil {
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

// readAddrSpec is used to read AddrSpec.
// Expects an address type byte, follwed by the address and port
func readAddrSpec(r io.Reader) (*AddrSpec, error) {
	d := &AddrSpec{}

	// Get the address type
	addrType := []byte{0}
	if _, err := r.Read(addrType); err != nil {
		return nil, err
	}

	fmt.Println(addrType)

	// Handle on a per type basis
	switch addrType[0] {
	case ipv4Address:
		addr := make([]byte, 4)
		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
			return nil, err
		}
		d.IP = net.IP(addr)

	case ipv6Address:
		addr := make([]byte, 16)
		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
			return nil, err
		}
		d.IP = net.IP(addr)

	case fqdnAddress:
		if _, err := r.Read(addrType); err != nil {
			return nil, err
		}
		fmt.Println(addrType)
		addrLen := int(addrType[0])
		fqdn := make([]byte, addrLen)
		if _, err := io.ReadAtLeast(r, fqdn, addrLen); err != nil {
			return nil, err
		}
		fmt.Println(fqdn)
		d.FQDN = string(fqdn)

	default:
		return nil, unrecognizedAddrType
	}

	// Read the port
	port := []byte{0, 0}
	if _, err := io.ReadAtLeast(r, port, 2); err != nil {
		return nil, err
	}

	fmt.Printf("port %+v\n", port)
	d.Port = (int(port[0]) << 8) | int(port[1])

	return d, nil
}

func sendReply(w io.Writer, resp uint8, addr *AddrSpec) error {
	// Format the address
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case addr == nil:
		addrType = ipv4Address
		addrBody = []byte{0, 0, 0, 0}
		addrPort = 0

	case addr.FQDN != "":
		addrType = fqdnAddress
		addrBody = append([]byte{byte(len(addr.FQDN))}, addr.FQDN...)
		addrPort = uint16(addr.Port)

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

	// Format the message
	msg := make([]byte, 6+len(addrBody))
	msg[0] = socks5Version
	msg[1] = resp
	msg[2] = 0 // Reserved
	msg[3] = addrType
	copy(msg[4:], addrBody)
	msg[4+len(addrBody)] = byte(addrPort >> 8)
	msg[4+len(addrBody)+1] = byte(addrPort & 0xff)

	fmt.Printf("<- %v\n", msg)

	// Send the message
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
