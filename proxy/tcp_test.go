package proxy

import (
	"net"
	"sync"
	"testing"
	"time"
    "errors"
    "bytes"
    "github.com/Einsteiniumeinsteinian/dns-over-tls-proxy/utility"

)

func TestStartTCPListener(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1) 

	// Start the listener in a goroutine
	go StartTCPListener(wg)

	time.Sleep(100 * time.Millisecond) 

	// Test: Connect to the server and ensure a response
	conn, err := net.Dial("tcp", ":"+utility.PortTCP)
	if err != nil {
		t.Fatalf("Failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	// Stop the listener
	wg.Done() // Signal the listener to shut down
	wg.Wait()  // Wait for it to finish
}


type mockConn struct {
    buf  bytes.Buffer
    fail bool
    data []byte
}
func (m *mockConn) Read(b []byte) (n int, err error) {
    if m.fail {
        return 0, errors.New("read error")
    }
    return copy(b, m.data), nil  // Simulate reading from the data slice
}

func (m *mockConn) Write(p []byte) (n int, err error) {
	if m.fail {
		return 0, errors.New("write error")
	}
	return m.buf.Write(p)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t time.Time) error {
    return nil // You might want to implement deadline logic if needed
}
func (m *mockConn) SetReadDeadline(t time.Time) error {
    return nil // You might want to implement deadline logic if needed
}
func (m *mockConn) SetWriteDeadline(t time.Time) error {
    return nil // You might want to implement deadline logic if needed
}
func TestSendErrorResponseTCP(t *testing.T) {
	// Create a mock connection
	conn := &mockConn{}

	// Call the function with an error
	errMsg := "some error"
	err := errors.New(errMsg)
	sendErrorResponseTCP(conn, err)

	// Check if the error message was written to the connection buffer
	expectedMsg := "DNS resolution failed: some error"
	if conn.buf.String() != expectedMsg {
		t.Errorf("Expected message %q, but got %q", expectedMsg, conn.buf.String())
	}
}
