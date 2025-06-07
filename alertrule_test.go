package librenms_test

import (
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"
	"github.com/stretchr/testify/require"
)

const (
	testEndpointAlertRules = "/api/v0/rules"
	testEndpointAlertRule  = "/api/v0/rules/1"
	testAlertRuleID        = 1
	testBuilderValue       = "{\"condition\":\"AND\",\"rules\":[{\"id\":\"macros.port_down\",\"field\":\"macros.port_down\",\"type\":\"integer\",\"input\":\"radio\",\"operator\":\"equal\",\"value\":\"1\"}],\"valid\":true}"
)

// This init function will register handlers for alert rule-related API endpoints.
func init() {
	handleEndpoint(testEndpointAlertRule, mockResponses{
		http.MethodDelete: loadMockResponse("delete_alertrule_200.json"),
		http.MethodGet:    loadMockResponse("get_alertrule_200.json"),
	})

	handleEndpoint(testEndpointAlertRules, mockResponses{
		http.MethodGet:  loadMockResponse("get_alertrules_200.json"),
		http.MethodPost: loadMockResponse("create_alertrule_200.json"),
		http.MethodPut:  loadMockResponse("update_alertrule_200.json"),
	})
}

func TestClient_GetAlertRule(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.GetAlertRule(testAlertRuleID)

	r.NoError(err, "GetAlertRule returned an error")
	r.NotNil(resp, "GetAlertRule response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
	r.Equal(1, resp.Count, "Expected count 1")
	r.Len(resp.Rules, 1, "Expected 1 alert rule")

	rule := resp.Rules[0]
	r.Equal(1, rule.ID, "Expected AlertRule ID 1")
	r.Equal("Device Down! Due to no ICMP response.", rule.Name, "Unexpected name")
}

func TestClient_GetAlertRules(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.GetAlertRules()

	r.NoError(err, "GetAlertRules returned an error")
	r.NotNil(resp, "GetAlertRules response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
	r.Equal(12, resp.Count, "Expected count 12")
	r.Len(resp.Rules, 12, "Expected 12 alert rules")

	rule := resp.Rules[0]
	r.Equal(1, rule.ID, "Expected AlertRule ID 1")
	r.Equal("Device Down! Due to no ICMP response.", rule.Name, "Unexpected Name")
}

func TestClient_CreateAlertRule(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the alert rule to create using AlertRuleCreateRequest
	newAlertRuleRequest := librenms.AlertRuleCreateRequest{
		Name:     "Test Alert Rule",
		Notes:    "This is a test alert rule",
		Devices:  []int{1, 2, 3},
		Disabled: false,
		Builder:  testBuilderValue,
	}

	createResp, err := testAPIClient.CreateAlertRule(&newAlertRuleRequest)

	r.NoError(err, "CreateAlertRule returned an error")
	r.NotNil(createResp, "CreateAlertRule response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
}

func TestClient_DeleteAlertRule(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.DeleteAlertRule(testAlertRuleID)

	r.NoError(err, "DeleteAlertRule returned an error")
	r.NotNil(resp, "DeleteAlertRule response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateAlertRule(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the alert rule to create using AlertRuleUpdateRequest
	payload := librenms.AlertRuleUpdateRequest{
		ID: testAlertRuleID,
		AlertRuleCreateRequest: librenms.AlertRuleCreateRequest{
			Name:     "Test Alert Rule",
			Notes:    "This is a test alert rule",
			Devices:  []int{1, 2, 3},
			Disabled: false,
			Builder:  testBuilderValue,
		},
	}

	resp, err := testAPIClient.UpdateAlertRule(&payload)

	r.NoError(err, "UpdateAlertRule returned an error")
	r.NotNil(resp, "UpdateAlertRule response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateAlertRule_EmptyDevices(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the alert rule to create using AlertRuleUpdateRequest
	payload := librenms.AlertRuleUpdateRequest{
		ID: testAlertRuleID,
		AlertRuleCreateRequest: librenms.AlertRuleCreateRequest{
			Name:     "Test Alert Rule",
			Notes:    "This is a test alert rule",
			Disabled: false,
			Builder:  testBuilderValue,
		},
	}

	resp, err := testAPIClient.UpdateAlertRule(&payload)

	r.NoError(err, "UpdateAlertRule returned an error")
	r.NotNil(resp, "UpdateAlertRule response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateAlertRule_MissingID(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	// Define the alert rule to create using AlertRuleUpdateRequest
	payload := librenms.AlertRuleUpdateRequest{
		AlertRuleCreateRequest: librenms.AlertRuleCreateRequest{
			Name:     "Test Alert Rule",
			Notes:    "This is a test alert rule",
			Devices:  []int{1, 2, 3},
			Disabled: false,
			Builder:  testBuilderValue,
		},
	}

	_, err := testAPIClient.UpdateAlertRule(&payload)

	r.Error(err, "UpdateAlertRule returned an error")
	r.Equal("rule ID is required for updating an alert rule", err.Error(), "Unexpected error message for missing ID")
}
