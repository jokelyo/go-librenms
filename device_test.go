package librenms_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	// Import the package under test
	"github.com/jokelyo/go-librenms"
	"github.com/stretchr/testify/require" // Changed from assert to require
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockGetDeviceResponse := loadMockResponse("get_device_200.json")
	mockCreateDeviceResponse := loadMockResponse("create_device_200.json")
	mockUpdateDeviceResponse := loadMockResponse("update_device_200.json")

	// Handler for /api/v0/devices/:id endpoint
	mux.HandleFunc("/api/v0/devices/1.1.1.1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(mockGetDeviceResponse)
			if err != nil {
				log.Printf("Error writing mockGetDeviceResponse: %v", err)
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		case http.MethodPatch:
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(mockUpdateDeviceResponse)
			if err != nil {
				log.Printf("Error writing mockUpdateDeviceResponse: %v", err)
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for /api/v0/devices.", r.Method), http.StatusMethodNotAllowed)
		}
	})

	// Handler for the /api/v0/devices endpoint
	mux.HandleFunc("/api/v0/devices/", func(w http.ResponseWriter, r *http.Request) { // Added trailing slash
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(mockCreateDeviceResponse)
			if err != nil {
				log.Printf("Error writing mockCreateDeviceResponse: %v", err)
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		default: // Catches GET and any other methods for /api/v0/devices
			http.Error(w, fmt.Sprintf("Method %s not implemented for /api/v0/devices.", r.Method), http.StatusMethodNotAllowed)
		}
	})
}

func TestClient_GetDevice(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	deviceResp, err := testAPIClient.GetDevice("1.1.1.1")

	r.NoError(err, "GetDevice returned an error")
	r.NotNil(deviceResp, "GetDevice response is nil")

	r.Equal("ok", deviceResp.Status, "Expected status 'ok'")
	r.Equal(1, deviceResp.Count, "Expected count 1")
	r.Len(deviceResp.Devices, 1, "Expected 1 device")

	device := deviceResp.Devices[0]
	r.Equal(1, device.DeviceID, "Expected DeviceID 1")
	r.Equal("1.1.1.1", device.Hostname, "Expected Hostname '1.1.1.1'")

	// verify a Bool field unmarshals correctly
	r.Equal(librenms.Bool(true), device.SNMPDisable, "Expected SNMPDisable true (1)")
}

func TestClient_CreateDevice(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the device to create using DeviceCreateRequest
	newDeviceRequest := librenms.DeviceCreateRequest{
		Hostname:      "192.168.10.5",
		SNMPVersion:   "v2c",    // Corresponds to 'snmpver' in JSON for request
		SNMPCommunity: "public", // Corresponds to 'community' in JSON for request
	}

	createResp, err := testAPIClient.CreateDevice(&newDeviceRequest)

	r.NoError(err, "CreateDevice returned an error")
	r.NotNil(createResp, "CreateDevice response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(1, createResp.Count, "Expected count 1")
	r.Len(createResp.Devices, 1, "Expected 1 device in response")

	// Verify the details of the device in the response
	// Note: The response device structure is librenms.Device
	deviceInResponse := createResp.Devices[0]
	r.Equal("192.168.10.5", deviceInResponse.Hostname, "Expected Hostname '192.168.10.5'")
	r.Equal("v2c", deviceInResponse.SNMPVersion, "Expected snmpver 'v2c'")

	// Verify that SNMPCommunity is not empty
	r.NotEmpty(deviceInResponse.Community, "Expected SNMPCommunity to be set in the response")
	r.Equal("public", *deviceInResponse.Community, "Expected SNMPCommunity 'public' in response")
}

func TestClient_UpdateDevice(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	payload := &librenms.DeviceUpdateRequest{
		Field: []string{"hardware", "port_association_mode"},
		Value: []any{"New Hardware", 2},
	}
	deviceResp, err := testAPIClient.UpdateDevice("1.1.1.1", payload)

	r.NoError(err, "GetDevice returned an error")
	r.NotNil(deviceResp, "GetDevice response is nil")

	r.Equal("ok", deviceResp.Status, "Expected status 'ok'")
	r.Equal("Device fields have been updated", deviceResp.Message, "Update message mismatch")
}
