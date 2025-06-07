package librenms_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jokelyo/go-librenms" // Import the package under test
	"github.com/stretchr/testify/require"
)

type (
	// mockResponses is a map of HTTP methods to their corresponding mock JSON responses,
	// which is used in the handleEndpoint() function.
	//
	// Use loadMockResponse() to load the JSON data from the fixtures directory.
	mockResponses map[string][]byte
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

// loadMockResponse is a helper function to load mock JSON responses from the fixtures directory.
func loadMockResponse(filename string) []byte {
	path := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read mock response file '%s': %v", path, err)
	}
	return data
}

// handleEndpoint is a helper function to register an endpoint with its methods and mock files.
//
// This function simplifies the process of setting up endpoints in tests, but if you
// need to customize the response handling (e.g., for different status codes or headers),
// you can define the handler function directly in the test file.
//
// For example:
//
//	handleEndpoint("/api/v0/alerts", mockResponses{
//	    http.MethodGet: loadMockResponse("get_alerts_200.json"),
//	})
func handleEndpoint(path string, methodsToMockResponses mockResponses) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodDelete, http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodPut:
			mockResponse, ok := methodsToMockResponses[r.Method]
			if !ok {
				notImplemented(path, w, r)
				return
			}
			_, err = w.Write(mockResponse)
		default:
			notImplemented(path, w, r)
			return
		}
		handleWriteErr(err, w)
	})
}

// handleWriteErr is a helper function to handle errors when writing responses in mux handlers.
func handleWriteErr(err error, w http.ResponseWriter) {
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// notImplemented is a helper function to handle unsupported HTTP methods for endpoints.
func notImplemented(path string, w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", path, r.Method), http.StatusMethodNotAllowed)
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
