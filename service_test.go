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
	testServiceID        = 2
	testEndpointServices = "/api/v0/services"
	testEndpointService  = "/api/v0/services/2"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockCreateServiceResponse := loadMockResponse("create_service_200.json")
	mockDeleteServiceResponse := loadMockResponse("delete_service_200.json")
	mockGetServicesResponse := loadMockResponse("get_services_200.json")
	mockGetServiceMembersResponse := loadMockResponse("get_service_200.json")
	mockUpdateServiceResponse := loadMockResponse("update_service_200.json")

	mux.HandleFunc(testEndpointService, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodDelete:
			_, err = w.Write(mockDeleteServiceResponse)
		case http.MethodGet:
			_, err = w.Write(mockGetServiceMembersResponse)
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

	mux.HandleFunc(testEndpointServices, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetServicesResponse)
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(mockCreateServiceResponse)
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
	r.Equal("GCP", service.Name, "Expected Service 'GCP'")
}

func TestClient_GetServicesForHost(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	serviceResp, err := testAPIClient.GetServicesForHost("4")

	r.NoError(err, "GetServicesForHost returned an error")
	r.NotNil(serviceResp, "GetServicesForHost response is nil")

	r.Equal("ok", serviceResp.Status, "Expected status 'ok'")
	// r.Equal(1, serviceResp.Count, "Expected count 1")
	r.Len(serviceResp.Services, 1, "Expected 1 services")

	service := serviceResp.Services[0]
	r.Equal(1, service.ID, "Expected Device ID 1")
}

func TestClient_CreateService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the service
	rules := librenms.ServiceRuleContainer{
		Condition: "AND",
		Rules: []librenms.ServiceRule{
			{
				ID:       "devices.sysDescr",
				Field:    "devices.sysDescr",
				Type:     "string",
				Input:    "text",
				Operator: "contains",
				Value:    "Linux",
			},
		},
		Joins: make([][]string, 0),
		Valid: true,
	}

	newServiceRequest := librenms.ServiceCreateRequest{
		Name:  "Test Service",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "dynamic",
	}

	createResp, err := testAPIClient.CreateService(&newServiceRequest)

	r.NoError(err, "CreateService returned an error")
	r.NotNil(createResp, "CreateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_CreateServiceNested(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the service. This definition makes no sense, but it doesn't matter.
	rules := librenms.ServiceRuleContainer{
		Condition: "AND",
		Rules: []librenms.ServiceRule{
			{
				Condition: "AND",
				Rules: []librenms.ServiceRule{
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
				},
			},
			{
				Condition: "AND",
				Rules: []librenms.ServiceRule{
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
					{
						ID:       "devices.sysDescr",
						Field:    "devices.sysDescr",
						Type:     "string",
						Input:    "text",
						Operator: "contains",
						Value:    "Linux",
					},
				},
			},
		},
		Joins: make([][]string, 0),
		Valid: true,
	}

	newServiceRequest := librenms.ServiceCreateRequest{
		Name:  "Test Service",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "dynamic",
	}

	createResp, err := testAPIClient.CreateService(&newServiceRequest)

	r.NoError(err, "CreateService returned an error")
	r.NotNil(createResp, "CreateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_CreateServiceStatic(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	newServiceRequest := librenms.ServiceCreateRequest{
		Name:    "Test Service",
		Devices: []int{1, 2},
		Type:    "static",
	}

	createResp, err := testAPIClient.CreateService(&newServiceRequest)

	r.NoError(err, "CreateService returned an error")
	r.NotNil(createResp, "CreateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_DeleteService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.DeleteService("4")

	r.NoError(err, "DeleteService returned an error")
	r.NotNil(resp, "DeleteService response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateService(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the service
	rules := librenms.ServiceRuleContainer{
		Condition: "AND",
		Rules: []librenms.ServiceRule{
			{
				ID:       "devices.sysDescr",
				Field:    "devices.sysDescr",
				Type:     "string",
				Input:    "text",
				Operator: "contains",
				Value:    "Windows",
			},
		},
		Joins: make([][]string, 0),
		Valid: true,
	}

	deviceServiceRequest := librenms.ServiceUpdateRequest{
		Name:  "Test Service",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "Dynamic",
	}

	createResp, err := testAPIClient.UpdateService("4", &deviceServiceRequest)

	r.NoError(err, "UpdateService returned an error")
	r.NotNil(createResp, "UpdateService response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
}
