package utility

import (
	"os"
	"fmt"
	"log"
	"github.com/joho/godotenv"
)

var (
    dnsServerHost     string
    dnsServerHostName string
    dnsServerPortTLS  string
	PortTCP			  string	
	PortUDP			  string
	UDP				  string
	TCP				  string
)

// init initializes environment variables by loading them from env

func init() {
    // Attempt to load .env file (ignore error if it doesn't exist)
    _ = godotenv.Load()

    // Load environment variables and handle errors
    var err error

    dnsServerHost, err = getEnvOrDefault("DNS_SERVER_HOST")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_HOST: %v", err)
    }

    dnsServerHostName, err = getEnvOrDefault("DNS_SERVER_HOST_NAME")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_HOST_NAME: %v", err)
    }

    dnsServerPortTLS, err = getEnvOrDefault("DNS_SERVER_PORT_TLS")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_PORT_TLS: %v", err)
    }

    PortTCP, err = getEnvOrDefault("PORT_TCP") 
    if err != nil {
        log.Fatalf("Error loading PORT_TCP: %v", err)
    }

    PortUDP, err = getEnvOrDefault("PORT_UDP") 
    if err != nil {
        log.Fatalf("Error loading PORT_TCP: %v", err)
    }

    UDP, _ = getEnvOrDefault("UDP") 

    TCP, _ = getEnvOrDefault("TCP") 

}

func getEnvOrDefault(key string) (string, error) {
    value, exists := os.LookupEnv(key)
    if !exists {
        return "", fmt.Errorf("environment variable %s not set", key) 
    }
    return value, nil
}
