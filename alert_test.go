package librenms_test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

const (
	testAlertID             = 9
	testEndpointAlerts      = "/api/v0/alerts"
	testEndpointAlert       = "/api/v0/alerts/9"
	testEndpointAlertUnmute = "/api/v0/alerts/unmute/9"
)

// This init function will register handlers for alert-related API endpoints.
func init() {
	handleEndpoint(testEndpointAlerts, mockResponses{
		http.MethodGet: loadMockResponse("get_alerts_200.json"),
	})

	handleEndpoint(testEndpointAlert, mockResponses{
		http.MethodGet: loadMockResponse("get_alert_200.json"),
		http.MethodPut: loadMockResponse("ack_alert_200.json"),
	})

	handleEndpoint(testEndpointAlertUnmute, mockResponses{
		http.MethodPut: loadMockResponse("unmute_alert_200.json"),
	})
}

func TestClient_AckAlert(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.AckAlert(
		testAlertID,
		&librenms.AlertAckRequest{UntilClear: true},
	)

	r.NoError(err, "AckAlert returned an error")
	r.NotNil(resp, "AckAlert response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
	r.Equal("Alert has been acknowledged", resp.Message, "Expected acknowledgment message")
}

func TestClient_AckAlert_NilPayload(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.AckAlert(testAlertID, nil)

	r.NoError(err, "AckAlert returned an error")
	r.NotNil(resp, "AckAlert response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
	r.Equal("Alert has been acknowledged", resp.Message, "Expected acknowledgment message")
}

func TestClient_GetAlert(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	alertResp, err := testAPIClient.GetAlert(testAlertID)

	r.NoError(err, "GetAlert returned an error")
	r.NotNil(alertResp, "GetAlert response is nil")

	r.Equal("ok", alertResp.Status, "Expected status 'ok'")
	r.Equal(1, alertResp.Count, "Expected count 1")
	r.Len(alertResp.Alerts, 1, "Expected 1 alerts")

	alert := alertResp.Alerts[0]
	r.Equal(testAlertID, alert.ID, "Expected AlertID "+strconv.Itoa(testAlertID))
	r.Equal(6, alert.DeviceID, "Expected DeviceID 6")
	r.Equal("warning", alert.Severity, "Expected Severity 'warning'")
}

func TestClient_GetAlerts(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	alertResp, err := testAPIClient.GetAlerts(
		librenms.NewAlertsQuery().SetState(1),
	)

	r.NoError(err, "GetAlerts returned an error")
	r.NotNil(alertResp, "GetAlerts response is nil")

	r.Equal("ok", alertResp.Status, "Expected status 'ok'")
	r.Equal(6, alertResp.Count, "Expected count 6")
	r.Len(alertResp.Alerts, 6, "Expected 6 alerts")

	alert := alertResp.Alerts[0]
	r.Equal(15, alert.ID, "Expected AlertID 15")
	r.Equal(3, alert.RuleID, "Expected Alert RuleID 3")
}

func TestClient_GetAlerts_NilPayload(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	alertResp, err := testAPIClient.GetAlerts(nil)

	r.NoError(err, "GetAlerts returned an error")
	r.NotNil(alertResp, "GetAlerts response is nil")

	r.Equal("ok", alertResp.Status, "Expected status 'ok'")
	r.Equal(6, alertResp.Count, "Expected count 6")
	r.Len(alertResp.Alerts, 6, "Expected 6 alerts")

	alert := alertResp.Alerts[0]
	r.Equal(15, alert.ID, "Expected AlertID 15")
	r.Equal(3, alert.RuleID, "Expected Alert RuleID 3")
}

func TestClient_UnmuteAlert(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.UnmuteAlert(testAlertID)

	r.NoError(err, "UnmuteAlert returned an error")
	r.NotNil(resp, "UnmuteAlert response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
	r.Equal("Alert has been unmuted", resp.Message, "Unexpected unmute message")
}
