package utility

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	maxPacketSize       = 1024                     // maximum DNS packet size
	tlsHandshakeTimeout = 10 * time.Second         // TLS handshake timeout
	tlsReadTimeout      = 5 * time.Second          // TLS read timeout
)

var ResolveDNSOverTLS = resolveDNSOverTLS

type TLSDialer func(network, addr string, config *tls.Config) (*tls.Conn, error)

// CustomError represents a custom error type with additional context.
type CustomError struct {
	Context string
	Err     error
}

// Error returns the error message for CustomError.
func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", ce.Context, ce.Err)
}

// verifyTLSCertificate checks that the certificate is valid and has not expired
func verifyTLSCertificate(conn *tls.Conn) error {
	log.Printf("Validating certificate for %s\n", dnsServerHostName)
	err := conn.VerifyHostname(dnsServerHostName)
	if err != nil {
		return &CustomError{Context: "failed to verify hostname", Err: err}
	}
	expires := conn.ConnectionState().PeerCertificates[0].NotAfter
	if expires.Before(time.Now()) {
		return &CustomError{Context: "certificate has expired", Err: errors.New("certificate has expired")}
	}

	return nil
}

func establishTLSConnection(dialer TLSDialer, config *tls.Config) (*tls.Conn, error) {
	return dialer("tcp", dnsServerHost+":"+dnsServerPortTLS, config)
}

func writeTLS(conn *tls.Conn, buf []byte) error {
	if _, err := conn.Write(buf); err != nil {
		return &CustomError{Context: "TLS write failed", Err: err}
	}
	return nil
}

func readTLS(conn *tls.Conn) ([]byte, error) {
	resp := make([]byte, maxPacketSize)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, &CustomError{Context: "TLS read failed", Err: err}
	}
	return resp[:n], nil
}

func setTLSDeadline(conn *tls.Conn, timeout time.Duration) error {
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return &CustomError{Context: "TLS deadline set failed", Err: err}
	}
	return nil
}

func resolveDNSOverTLS(buf []byte, dialer TLSDialer, config *tls.Config) ([]byte, error) {
	conn, err := establishTLSConnection(dialer, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := verifyTLSCertificate(conn); err != nil {
		return nil, err
	}

	if err := setTLSDeadline(conn, tlsHandshakeTimeout); err != nil {
		return nil, err
	}
	if err := writeTLS(conn, buf); err != nil {
		return nil, err
	}

	if err := setTLSDeadline(conn, tlsReadTimeout); err != nil {
		return nil, err
	}
	resp, err := readTLS(conn)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
