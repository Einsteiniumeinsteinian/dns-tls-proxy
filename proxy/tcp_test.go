package proxy

import (
	"net"
	"sync"
	"testing"
	"time"
    "errors"
    "bytes"

)

func TestStartTCPListener(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1) 

	// Start the listener in a goroutine
	go StartTCPListener(wg)

	// Give it a moment to start up
	time.Sleep(100 * time.Millisecond) 

	// Test: Connect to the server and ensure a response
	conn, err := net.Dial("tcp", ":"+PortTCP)
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

// package proxy

// import (
//     "io"
//     "net"
//     "testing"
//     "time"
//     "sync"
// )

// // MockConn is a mock implementation of net.Conn for testing.
// type MockConn struct {
//     Closed bool
// }

// func (mc *MockConn) Read(b []byte) (n int, err error) {
//     return 0, io.EOF // Mocking a closed connection
// }

// func (mc *MockConn) Write(b []byte) (n int, err error) {
//     return len(b), nil
// }

// func (mc *MockConn) Close() error {
//     mc.Closed = true
//     return nil
// }

// func (mc *MockConn) LocalAddr() net.Addr {
//     return nil
// }

// func (mc *MockConn) RemoteAddr() net.Addr {
//     return nil
// }

// func (mc *MockConn) SetDeadline(t time.Time) error {
//     return nil
// }

// func (mc *MockConn) SetReadDeadline(t time.Time) error {
//     return nil
// }

// func (mc *MockConn) SetWriteDeadline(t time.Time) error {
//     return nil
// }

// // TestStartTCPListener tests the StartTCPListener function.
// func TestStartTCPListener(t *testing.T) {
//     // Create a wait group for synchronization
//     var wg sync.WaitGroup

//     defer wg.Done()
//     wg.Add(1)

//     // Start the TCP listener in a goroutine
//     go StartTCPListener(&wg)

//     // Create a mock TCP connection
//     mockConn := &MockConn{}

//     // Simulate handling TCP connection with the mock connection
//     go HandleTCPConnection(mockConn)

//     // Wait for the listener to finish
//     wg.Wait()

//     // Assert that the connection is closed
//     if !mockConn.Closed {
//         t.Error("Expected connection to be closed")
//     }
// }

// MockDNSResolver is a mock implementation of the DNSResolver interface
// type MockDNSResolver struct{}

// // ResolveDNSOverTLS returns a mock DNS response
// func (r *MockDNSResolver) ResolveDNSOverTLS(query []byte, dial func(network, addr string, config *tls.Config) (*tls.Conn, error), config *tls.Config) ([]byte, error) {
//     return []byte("Mock DNS response"), nil
// }

// func TestStartTCPListener(t *testing.T) {
//     // Create a mock DNS resolver
//     mockResolver := &MockDNSResolver{}

//     wg := &sync.WaitGroup{}
//     wg.Add(1)
//     // Start the TCP listener with the mock resolver
//     go StartTCPListener(wg, mockResolver)

//     // Wait for the server to start (you might need to adjust the timeout)
//     time.Sleep(100 * time.Millisecond)

//     // Simulate a client connection
//     conn, err := net.Dial("tcp", "localhost:"+portTCP)
//     if err != nil {
//         t.Fatalf("Failed to connect to server: %v", err)
//     }
//     defer conn.Close()

//     // Send a query
//     query := []byte("example.com")
//     _, err = conn.Write(query)
//     if err != nil {
//         t.Fatalf("Failed to send query: %v", err)
//     }

//     // Receive the response
//     buf := make([]byte, 1024)
//     n, err := conn.Read(buf)
//     if err != nil && err != io.EOF {
//         t.Fatalf("Failed to read response: %v", err)
//     }
//     response := buf[:n]

//     // Assert the response
//     expectedResponse := []byte("Mock DNS response")
//     if !bytes.Equal(response, expectedResponse) {
//         t.Errorf("Unexpected response: got %q, want %q", response, expectedResponse)
//     }

//     // Stop the server gracefully
//     conn.Close()
//     wg.Wait()
// }

// func TestHandleTCPConnectionError(t *testing.T) {
//     // Create a mock DNS resolver that returns an error
//     mockResolver := &MockDNSResolver{}
    
//     // Simulate a connection
//     server, client := net.Pipe()
//     defer server.Close()

//     // Start handling the connection
//     go proxy.HandleTCPConnection(server, mockResolver) 

//     // Send some data
//     client.Write([]byte("test query"))

//     // Expect an error message
//     expectedError := []byte("DNS resolution failed: ") // Adjust the expected error based on your mockResolver implementation
//     buf := make([]byte, len(expectedError))
//     if _, err := client.Read(buf); err != nil {
//         t.Fatalf("Failed to read error message: %v", err)
//     }

//     // Assert the error message is correct
//     if !bytes.Equal(buf, expectedError) {
//         t.Errorf("Unexpected error message: got %q, want %q", buf, expectedError)
//     }
// }
