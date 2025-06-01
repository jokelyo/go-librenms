package librenms_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

const (
	testEndpointDeviceGroups = "/api/v0/devicegroups"
	testEndpointDeviceGroup  = "/api/v0/devicegroups/4"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockCreateDeviceGroupResponse := loadMockResponse("create_devicegroup_201.json")
	mockGetDeviceGroupsResponse := loadMockResponse("get_devicegroups_200.json")
	mockUpdateDeviceGroupResponse := loadMockResponse("update_devicegroup_200.json")

	// Handler for /api/v0/devices/:id endpoint
	mux.HandleFunc(testEndpointDeviceGroup, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPatch:
			_, err = w.Write(mockUpdateDeviceGroupResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointDeviceGroup, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Handler for the /api/v0/devices endpoint
	mux.HandleFunc(testEndpointDeviceGroups, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetDeviceGroupsResponse)
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(mockCreateDeviceGroupResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointDeviceGroups, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func TestClient_GetDeviceGroups(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	groupResp, err := testAPIClient.GetDeviceGroups()

	r.NoError(err, "GetDeviceGroups returned an error")
	r.NotNil(groupResp, "GetDeviceGroups response is nil")

	r.Equal("ok", groupResp.Status, "Expected status 'ok'")
	r.Equal(3, groupResp.Count, "Expected count 3")
	r.Len(groupResp.Groups, 3, "Expected 3 device groups")

	group := groupResp.Groups[0]
	r.Equal(1, group.ID, "Expected GroupID 1")
	r.Equal("GCP", group.Name, "Expected Group 'GCP'")

}

func TestClient_CreateDeviceGroup(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the device group
	rules := librenms.DeviceGroupRuleContainer{
		Condition: "AND",
		Rules: []librenms.DeviceGroupRule{
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

	newDeviceGroupRequest := librenms.DeviceGroupCreateRequest{
		Name:  "Test Group",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "Dynamic",
	}

	createResp, err := testAPIClient.CreateDeviceGroup(&newDeviceGroupRequest)

	r.NoError(err, "CreateDeviceGroup returned an error")
	r.NotNil(createResp, "CreateDeviceGroup response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_UpdateDeviceGroup(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the device group
	rules := librenms.DeviceGroupRuleContainer{
		Condition: "AND",
		Rules: []librenms.DeviceGroupRule{
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

	deviceGroupRequest := librenms.DeviceGroupUpdateRequest{
		Name:  "Test Group",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "Dynamic",
	}

	createResp, err := testAPIClient.UpdateDeviceGroup("4", &deviceGroupRequest)

	r.NoError(err, "UpdateDeviceGroup returned an error")
	r.NotNil(createResp, "UpdateDeviceGroup response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
}
