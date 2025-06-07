package librenms_test

import (
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

// This init function will register handlers for service-related API endpoints.
func init() {
	handleEndpoint(testEndpointService, mockResponses{
		http.MethodDelete: loadMockResponse("delete_service_200.json"),
		http.MethodPatch:  loadMockResponse("update_service_200.json"),
	})

	handleEndpoint(testEndpointServiceDeviceID, mockResponses{
		http.MethodGet:  loadMockResponse("get_service_200.json"),
		http.MethodPost: loadMockResponse("create_service_200.json"),
	})

	handleEndpoint(testEndpointServices, mockResponses{
		http.MethodGet: loadMockResponse("get_services_200.json"),
	})
}

func TestClient_GetService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	serviceResp, err := testAPIClient.GetService(testServiceID)

	r.NoError(err, "GetService returned an error")
	r.NotNil(serviceResp, "GetService response is nil")

	r.Equal("ok", serviceResp.Status, "Expected status 'ok'")
	r.Equal(1, serviceResp.Count, "Expected count 1")
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
	r.Equal(3, serviceResp.Count, "Expected count 3")
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
	r.Equal(2, serviceResp.Count, "Expected count 2")
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

	resp, err := testAPIClient.UpdateService(
		testServiceID,
		librenms.NewServiceUpdateRequest().SetName("Fancy Test Service"),
	)

	r.NoError(err, "UpdateService returned an error")
	r.NotNil(resp, "UpdateService response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}
