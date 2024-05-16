package proxy

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"sync"

	"github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/utility"
)

var (
	portUDP       = os.Getenv("PORT_UDP")
	maxPacketSize = 1024 // maximum DNS packet size
)

func StartUDPListener(wg *sync.WaitGroup) {
	defer wg.Done()
	pc, err := net.ListenPacket("udp", ":"+portUDP)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v\n", err)
	}
	defer pc.Close()
	log.Printf("UDP server listening on port %s...\n", portUDP)
	for {
		buf := make([]byte, maxPacketSize)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Printf("UDP server error: %v\n", err)
			continue
		}
		go handleUDPConnection(pc, buf[:n], addr)
	}
}

func handleUDPConnection(pc net.PacketConn, buf []byte, addr net.Addr) {
	log.Printf("Received UDP packet\n")
	newBuf := []byte{0, byte(len(buf))}
	newBuf = append(newBuf, buf...)

	// Handle potential error from resolveTLS
	resp, err := utility.ResolveDNSOverTLS(newBuf, tls.Dial, &tls.Config{}) // Use tls.Dial and &tls.Config{}
	if err != nil {
		log.Printf("Failed to request domain resolution: %s\n", err)
		sendErrorResponseUDP(pc, addr, err)
		return
	}

	pc.WriteTo(resp[2:], addr)
	log.Printf("Server response sent\n")
}

// sendErrorResponseUDP sends an error response over UDP connection
func sendErrorResponseUDP(pc net.PacketConn, addr net.Addr, err error) {
	log.Printf("Sending error response: %v\n", err)
	errorMsg := []byte("DNS resolution failed: " + err.Error()) // Creating a simple error message
	_, writeErr := pc.WriteTo(errorMsg, addr)
	if writeErr != nil {
		log.Printf("Error sending error response: %v\n", writeErr)
	}
}
