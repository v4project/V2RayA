package tls

import (
	"crypto/tls"
	"fmt"
	"github.com/v2rayA/v2rayA/plugin/infra"
	"net"
	"net/url"
	"strings"
)

// Tls is a base Tls struct
type Tls struct {
	dialer     infra.Dialer
	addr       string
	serverName string
	skipVerify bool
	tlsConfig  *tls.Config
}

// NewTls returns a Tls infra.
func NewTls(s string, d infra.Dialer) (*Tls, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, newError("[Tls]").Base(err)
	}

	t := &Tls{
		dialer: d,
		addr:   u.Host,
	}

	query := u.Query()
	t.serverName = query.Get("host")
	if t.serverName == "" {
		colonPos := strings.LastIndex(t.addr, ":")
		if colonPos == -1 {
			colonPos = len(t.addr)
		}
		t.serverName = t.addr[:colonPos]
	}

	// skipVerify
	if query.Get("skipVerify") == "true" || query.Get("skipVerify") == "1" {
		t.skipVerify = true
	}

	t.tlsConfig = &tls.Config{
		ServerName:         t.serverName,
		InsecureSkipVerify: t.skipVerify,
	}

	return t, nil
}

// Addr returns forwarder's address.
func (s *Tls) Addr() string {
	if s.addr == "" {
		return s.dialer.Addr()
	}
	return s.addr
}

// Dial connects to the address addr on the network net via the infra.
func (s *Tls) Dial(network, addr string) (net.Conn, error) {
	return s.dial(network, addr)
}

func (s *Tls) dial(network, addr string) (conn net.Conn, err error) {
	rc, err := s.dialer.Dial("tcp", addr)
	if err != nil {
		return nil, newError(fmt.Sprintf("[Tls]: dial to %s", s.addr)).Base(err)
	}

	tlsConn := tls.Client(rc, s.tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}
	return tlsConn, err
}

// DialUDP connects to the given address via the infra.
func (s *Tls) DialUDP(network, addr string) (net.PacketConn, net.Addr, error) {
	//TODO
	return nil, nil, nil
}
