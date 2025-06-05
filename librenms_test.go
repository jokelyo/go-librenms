package librenms_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jokelyo/go-librenms" // Import the package under test
	"github.com/stretchr/testify/require"
)

var (
	testServer    *httptest.Server
	testAPIClient *librenms.Client
	// Global Mux for the test server. Handlers will be added by init() functions in test files.
	mux = http.NewServeMux()
)

func TestMain(m *testing.M) {
	// The mux is now global and populated by init() functions in test files.
	testServer = httptest.NewServer(mux)

	var err error // Declare err for librenms.New
	testAPIClient, err = librenms.New(testServer.URL+"/", "test-token-global")
	if err != nil {
		testServer.Close() // Clean up server if client creation fails
		log.Fatalf("Failed to create global test API client: %v", err)
	}

	// Run tests
	code := m.Run()

	// Teardown
	testServer.Close()
	os.Exit(code)
}

// Helper function to load mock JSON responses from the fixtures directory.
func loadMockResponse(filename string) []byte {
	path := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read mock response file '%s': %v", path, err)
	}
	return data
}

func TestClient_InvalidHostProtocol(t *testing.T) {
	r := require.New(t)

	// Test creating a client URL missing the protocol
	_, err := librenms.New("localhost:43433/", "test-token")
	r.Error(err, "Expected error when creating client with missing baseURL protocol")
	r.ErrorContains(err, "invalid base URL format", "Expected invalid base URL format error")
}

func TestClient_InvalidHostURI(t *testing.T) {
	r := require.New(t)

	// Test creating a client with an invalid trailing URI
	_, err := librenms.New("http://localhost:48325/api", "test-token")
	r.Error(err, "Expected error when creating client with invalid baseURL URI")
	r.ErrorContains(err, "invalid base URL format", "Expected invalid base URL format error")
}

func TestClient_ConnectionRefused(t *testing.T) {
	r := require.New(t)

	// Test creating a client with an unresponsive host
	client, err := librenms.New("http://localhost:48325/", "test-token")
	r.NoError(err, "Expected no error when creating client with unresponsive host")

	_, err = client.GetDevices(nil)
	r.Error(err, "Expected error when using client with unresponsive host")
	r.ErrorContains(err, "connection refused", "Expected connection refused error")
}
