package proxy

import (
	"crypto/tls"
	"io"
	"os"
	"log"
	"net"
	"sync"
	"github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/utility"
)

var (
	PortTCP             = os.Getenv("PORT_TCP")
)

func StartTCPListener(wg *sync.WaitGroup) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":"+PortTCP)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v\n", err)
	}
	defer ln.Close()
	log.Printf("TCP server listening on port %s...\n", PortTCP)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("TCP server error: %v\n", err)
			continue
		}
		go HandleTCPConnection(conn)
	}
}


func HandleTCPConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, maxPacketSize)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("TCP connection read error: %v\n", err)
		return
	}

	// Handle potential error from resolveTLS
	resp, err := utility.ResolveDNSOverTLS(buf[:n], tls.Dial, &tls.Config{}) // Use tls.Dial and &tls.Config{}
	if err != nil {
		log.Printf("Failed to resolve DNS over TLS: %v\n", err)
		sendErrorResponseTCP(conn, err)
		return
	}

	if _, err := conn.Write(resp); err != nil {
		log.Printf("TCP connection write error: %v\n", err)
		return
	}
}


// sendErrorResponse sends an error response over TCP connection
func sendErrorResponseTCP(conn net.Conn, err error) {
	log.Printf("Sending error response: %v\n", err)
	errorMsg := []byte("DNS resolution failed: " + err.Error()) // Creating a simple error message
	_, writeErr := conn.Write(errorMsg)
	if writeErr != nil {
		log.Printf("Error sending error response: %v\n", writeErr)
	}
}