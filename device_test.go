package librenms_test

import (
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"
	"github.com/stretchr/testify/require"
)

const (
	testEndpointDevices      = "/api/v0/devices"
	testEndpointDevicesSlash = "/api/v0/devices/"
	testEndpointDevice       = "/api/v0/devices/1.1.1.1"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	handleEndpoint(testEndpointDevice, mockResponses{
		http.MethodDelete: loadMockResponse("delete_device_200.json"),
		http.MethodGet:    loadMockResponse("get_device_200.json"),
		http.MethodPatch:  loadMockResponse("update_device_200.json"),
	})

	handleEndpoint(testEndpointDevicesSlash, mockResponses{
		http.MethodPost: loadMockResponse("create_device_200.json"),
	})

	handleEndpoint(testEndpointDevices, mockResponses{
		http.MethodGet: loadMockResponse("get_devices_200.json"),
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

func TestClient_GetDevices(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	deviceResp, err := testAPIClient.GetDevices(nil)

	r.NoError(err, "GetDevices returned an error")
	r.NotNil(deviceResp, "GetDevices response is nil")

	r.Equal("ok", deviceResp.Status, "Expected status 'ok'")
	r.Equal(3, deviceResp.Count, "Expected count 3")
	r.Len(deviceResp.Devices, 3, "Expected 3 devices")

	device := deviceResp.Devices[0]
	r.Equal(1, device.DeviceID, "Expected DeviceID 1")
	r.Equal("1.1.1.1", device.Hostname, "Expected Hostname '1.1.1.1'")

	// verify a Bool field unmarshals correctly
	r.Equal(librenms.Bool(true), device.SNMPDisable, "Expected SNMPDisable true (1)")

	// verify a Float64 field unmarshals correctly
	device = deviceResp.Devices[2]
	r.Equal(2, device.DeviceID, "Expected DeviceID 2")
	r.NotNil(device.Latitude, "Expected Latitude to be non-nil")
	r.Equal(-45.08624620, float64(*device.Latitude), "Expected Latitude -45.0862462")
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

func TestClient_DeleteDevice(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	deviceResp, err := testAPIClient.DeleteDevice("1.1.1.1")

	r.NoError(err, "DeleteDevice returned an error")
	r.NotNil(deviceResp, "DeleteDevice response is nil")

	r.Equal("ok", deviceResp.Status, "Expected status 'ok'")
	r.Equal(1, deviceResp.Count, "Expected count 1")
	r.Len(deviceResp.Devices, 1, "Expected 1 device")

	device := deviceResp.Devices[0]
	r.Equal(1, device.DeviceID, "Expected DeviceID 1")
	r.Equal("1.1.1.1", device.Hostname, "Expected Hostname '1.1.1.1'")
}

func TestClient_UpdateDevice(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	payload := &librenms.DeviceUpdateRequest{
		Field: []string{"hardware", "port_association_mode"},
		Data:  []any{"New Hardware", 2},
	}
	deviceResp, err := testAPIClient.UpdateDevice("1.1.1.1", payload)

	r.NoError(err, "GetDevice returned an error")
	r.NotNil(deviceResp, "GetDevice response is nil")

	r.Equal("ok", deviceResp.Status, "Expected status 'ok'")
	r.Equal("Device fields have been updated", deviceResp.Message, "Update message mismatch")
}
