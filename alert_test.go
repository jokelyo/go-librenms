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
	testAlertID             = 9
	testEndpointAlerts      = "/api/v0/alerts"
	testEndpointAlert       = "/api/v0/alerts/9"
	testEndpointAlertUnmute = "/api/v0/alerts/unmute/9"
)

// This init function will register handlers for device-related API endpoints.
func init() {
	mockAckAlertResponse := loadMockResponse("ack_alert_200.json")
	mockGetAlertsResponse := loadMockResponse("get_alerts_200.json")
	mockGetAlertResponse := loadMockResponse("get_alert_200.json")
	mockUnmuteAlertResponse := loadMockResponse("unmute_alert_200.json")

	mux.HandleFunc(testEndpointAlerts, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetAlertsResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointAlerts, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc(testEndpointAlert, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetAlertResponse)
		case http.MethodPut:
			_, err = w.Write(mockAckAlertResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointAlert, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc(testEndpointAlertUnmute, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			_, err = w.Write(mockUnmuteAlertResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", testEndpointAlert, r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
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
