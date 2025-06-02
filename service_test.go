package librenms_test

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

const (
	testServiceDeviceID         = "13"
	testServiceID               = 2
	testEndpointServices        = "/api/v0/services"
	testEndpointService         = "/api/v0/services/2"
	testEndpointServiceDeviceID = "/api/v0/services/13"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockCreateServiceResponse := loadMockResponse("create_service_200.json")
	mockDeleteServiceResponse := loadMockResponse("delete_service_200.json")
	mockGetServicesResponse := loadMockResponse("get_services_200.json")
	mockGetServiceMembersResponse := loadMockResponse("get_service_200.json")
	mockUpdateServiceResponse := loadMockResponse("update_service_200.json")

	// PATCH and DELETE for services/:id expect a service ID.
	mux.HandleFunc(testEndpointService, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodDelete:
			_, err = w.Write(mockDeleteServiceResponse)
		case http.MethodPatch:
			_, err = w.Write(mockUpdateServiceResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointService, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// POST and GET for services/:id expect a device ID or hostname.
	mux.HandleFunc(testEndpointServiceDeviceID, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetServiceMembersResponse)
		case http.MethodPost:
			_, err = w.Write(mockCreateServiceResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointServiceDeviceID, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc(testEndpointServices, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetServicesResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointServices, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func TestClient_GetService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	serviceResp, err := testAPIClient.GetService(testServiceID)

	r.NoError(err, "GetService returned an error")
	r.NotNil(serviceResp, "GetService response is nil")

	r.Equal("ok", serviceResp.Status, "Expected status 'ok'")
	// r.Equal(1, serviceResp.Count, "Expected count 1")
	r.Len(serviceResp.Services, 1, "Expected 1 services")

	service := serviceResp.Services[0]
	r.Equal(testServiceID, service.ID, "Expected ServiceID "+strconv.Itoa(testServiceID))
	r.Equal("check other thing", service.Name, "Unxpected service name")
}

func TestClient_GetServices(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	serviceResp, err := testAPIClient.GetServices()

	r.NoError(err, "GetServices returned an error")
	r.NotNil(serviceResp, "GetServices response is nil")

	r.Equal("ok", serviceResp.Status, "Expected status 'ok'")
	// r.Equal(3, serviceResp.Count, "Expected count 3") // seems like count is always 1, in v25.5
	r.Len(serviceResp.Services, 3, "Expected 3 services")

	service := serviceResp.Services[0]
	r.Equal(1, service.ID, "Expected ServiceID 1")
	r.Equal("check https cert", service.Name, "Expected Service name 'check https cert'")
}

func TestClient_GetServicesForHost(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	serviceResp, err := testAPIClient.GetServicesForHost(testServiceDeviceID)

	r.NoError(err, "GetServicesForHost returned an error")
	r.NotNil(serviceResp, "GetServicesForHost response is nil")

	r.Equal("ok", serviceResp.Status, "Expected status 'ok'")
	// r.Equal(1, serviceResp.Count, "Expected count 1")
	r.Len(serviceResp.Services, 2, "Expected 2 services")

	service := serviceResp.Services[0]
	r.Equal(1, service.ID, "Expected Service ID 1")
}

func TestClient_CreateService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	newServiceRequest := librenms.ServiceCreateRequest{
		Name:        "Test Service",
		Description: "This is a test service",
		IP:          "192.168.1.1",
		Param:       "-t 10 -c 5",
		Type:        "ping",
	}

	createResp, err := testAPIClient.CreateService(testServiceDeviceID, &newServiceRequest)

	r.NoError(err, "CreateService returned an error")
	r.NotNil(createResp, "CreateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal("Service ping has been added to device 2 (#1)", createResp.Message, "Unexpected message")
}

func TestClient_DeleteService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.DeleteService(testServiceID)

	r.NoError(err, "DeleteService returned an error")
	r.NotNil(resp, "DeleteService response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	deviceServiceRequest := librenms.ServiceUpdateRequest{
		Name: "Fancy Test Service",
	}

	createResp, err := testAPIClient.UpdateService(testServiceID, &deviceServiceRequest)

	r.NoError(err, "UpdateService returned an error")
	r.NotNil(createResp, "UpdateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
}
