package librenms_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockCreateDeviceGroupResponse := loadMockResponse("create_devicegroup_201.json")
	mockGetDeviceGroupsResponse := loadMockResponse("get_devicegroups_200.json")

	// Handler for /api/v0/devices/:id endpoint
	//mux.HandleFunc("/api/v0/devicegroups/1", func(w http.ResponseWriter, r *http.Request) {
	//	var err error
	//	w.Header().Set("Content-Type", "application/json")
	//	switch r.Method {
	//	case http.MethodDelete:
	//		_, err = w.Write(mockDeleteDeviceResponse)
	//	case http.MethodGet:
	//		_, err = w.Write(mockGetDeviceResponse)
	//	case http.MethodPatch:
	//		_, err = w.Write(mockUpdateDeviceResponse)
	//	default:
	//		http.Error(w, fmt.Sprintf("Method %s not implemented for /api/v0/devices.", r.Method), http.StatusMethodNotAllowed)
	//		return
	//	}
	//	if err != nil {
	//		log.Printf("Error writing response: %v", err)
	//		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	//	}
	//})

	// Handler for the /api/v0/devices endpoint
	mux.HandleFunc("/api/v0/devicegroups", func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetDeviceGroupsResponse)
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(mockCreateDeviceGroupResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for /api/v0/devicegroups.", r.Method), http.StatusMethodNotAllowed)
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
