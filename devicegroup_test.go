package librenms_test

import (
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

const (
	testEndpointDeviceGroups = "/api/v0/devicegroups"
	testEndpointDeviceGroup  = "/api/v0/devicegroups/4"
)

// This init function will register handlers for devicegroup-related API endpoints.
func init() {
	handleEndpoint(testEndpointDeviceGroup, mockResponses{
		http.MethodDelete: loadMockResponse("delete_devicegroup_200.json"),
		http.MethodGet:    loadMockResponse("get_devicegroup_200.json"),
		http.MethodPatch:  loadMockResponse("update_devicegroup_200.json"),
	})

	// Registering this endpoint outside of handleEndpoint() to mock the HTTP 201 POST response.
	mux.HandleFunc(testEndpointDeviceGroups, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(loadMockResponse("get_devicegroups_200.json"))
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(loadMockResponse("create_devicegroup_201.json"))
		default:
			notImplemented(testEndpointDeviceGroups, w, r)
			return
		}
		handleWriteErr(err, w)
	})
}

func TestClient_GetDeviceGroup(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	groupResp, err := testAPIClient.GetDeviceGroup("4")

	r.NoError(err, "GetDeviceGroup returned an error")
	r.NotNil(groupResp, "GetDeviceGroup response is nil")

	r.Equal("ok", groupResp.Status, "Expected status 'ok'")
	r.Equal(1, groupResp.Count, "Expected count 1")
	r.Len(groupResp.Groups, 1, "Expected 1 device groups")

	group := groupResp.Groups[0]
	r.Equal(4, group.ID, "Expected GroupID 4")
	r.Equal("NestedRules", group.Name, "Expected Group 'NestedRules'")
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

func TestClient_GetDeviceGroupMembers(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	groupResp, err := testAPIClient.GetDeviceGroupMembers("4")

	r.NoError(err, "GetDeviceGroupMembers returned an error")
	r.NotNil(groupResp, "GetDeviceGroupMembers response is nil")

	r.Equal("ok", groupResp.Status, "Expected status 'ok'")
	r.Equal(1, groupResp.Count, "Expected count 1")
	r.Len(groupResp.Devices, 1, "Expected 1 device groups")

	member := groupResp.Devices[0]
	r.Equal(6, member.ID, "Expected Device ID 6")
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
		Type:  "dynamic",
	}

	createResp, err := testAPIClient.CreateDeviceGroup(&newDeviceGroupRequest)

	r.NoError(err, "CreateDeviceGroup returned an error")
	r.NotNil(createResp, "CreateDeviceGroup response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_CreateDeviceGroupNested(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the rules for the device group. This definition makes no sense, but it doesn't matter.
	rules := librenms.DeviceGroupRuleContainer{
		Condition: "AND",
		Rules: []librenms.DeviceGroupRule{
			{
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
				Rules: []librenms.DeviceGroupRule{
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

	newDeviceGroupRequest := librenms.DeviceGroupCreateRequest{
		Name:  "Test Group",
		Rules: func() *string { s := rules.MustJSON(); return &s }(),
		Type:  "dynamic",
	}

	createResp, err := testAPIClient.CreateDeviceGroup(&newDeviceGroupRequest)

	r.NoError(err, "CreateDeviceGroup returned an error")
	r.NotNil(createResp, "CreateDeviceGroup response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_CreateDeviceGroupStatic(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	newDeviceGroupRequest := librenms.DeviceGroupCreateRequest{
		Name:    "Test Group",
		Devices: []int{1, 2},
		Type:    "static",
	}

	createResp, err := testAPIClient.CreateDeviceGroup(&newDeviceGroupRequest)

	r.NoError(err, "CreateDeviceGroup returned an error")
	r.NotNil(createResp, "CreateDeviceGroup response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
	r.Equal(4, createResp.ID, "Expected ID 4")
}

func TestClient_DeleteDeviceGroup(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.DeleteDeviceGroup("4")

	r.NoError(err, "DeleteDeviceGroup returned an error")
	r.NotNil(resp, "DeleteDeviceGroup response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
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
