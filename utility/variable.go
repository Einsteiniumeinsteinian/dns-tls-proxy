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

    dnsServerHost, err = GetEnvOrDefault("DNS_SERVER_HOST")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_HOST: %v", err)
    }

    dnsServerHostName, err = GetEnvOrDefault("DNS_SERVER_HOST_NAME")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_HOST_NAME: %v", err)
    }

    dnsServerPortTLS, err = GetEnvOrDefault("DNS_SERVER_PORT_TLS")
    if err != nil {
        log.Fatalf("Error loading DNS_SERVER_PORT_TLS: %v", err)
    }

    PortTCP, err = GetEnvOrDefault("PORT_TCP") 
    if err != nil {
        log.Fatalf("Error loading PORT_TCP: %v", err)
    }

    PortUDP, err = GetEnvOrDefault("PORT_UDP") 
    if err != nil {
        log.Fatalf("Error loading PORT_TCP: %v", err)
    }

    UDP, _ = GetEnvOrDefault("UDP") 

    TCP, _ = GetEnvOrDefault("TCP") 

}

func GetEnvOrDefault(key string) (string, error) {
    value, exists := os.LookupEnv(key)
    if !exists {
        return "", fmt.Errorf("environment variable %s not set", key) 
    }
    return value, nil
}
