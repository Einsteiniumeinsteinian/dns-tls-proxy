package main

import (
	"os"
	"sync"
	"github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/proxy"
	"log"
)


func main() {
    var wg sync.WaitGroup
    wg.Add(2)  // Change to wg.Add(0)

    // Get environment variables (or use command-line flags if preferred)
    startTCP := os.Getenv("TCP")
    startUDP := os.Getenv("UDP")

    if startTCP == "true" {
        wg.Add(1)
        go proxy.StartTCPListener(&wg)
    }

    if startUDP == "true" {
        wg.Add(1)
        go proxy.StartUDPListener(&wg)
    }

    if startTCP != "true" && startUDP != "true" {
        log.Fatal("Error: Either START_TCP or START_UDP must be set to 'true'.")
    }

    wg.Wait()
}