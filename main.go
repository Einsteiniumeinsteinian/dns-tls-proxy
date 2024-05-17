package main

import (
	"sync"
	"github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/proxy"
	"github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/utility"
	"log"
)


func main() {
    var wg sync.WaitGroup
    wg.Add(2)  // Change to wg.Add(0)

    if utility.TCP == "true" {
        wg.Add(1)
        go proxy.StartTCPListener(&wg)
    }

    if utility.UDP == "true" {
        wg.Add(1)
        go proxy.StartUDPListener(&wg)
    }

    if utility.TCP != "true" && utility.UDP != "true" {
        log.Fatal("Error: Either START_TCP or START_UDP must be set to 'true'.")
    }

    wg.Wait()
}
