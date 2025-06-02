package librenms

import (
	"fmt"
	"net/http"
)

const (
	alertRuleEndpoint = "rules"
)

type (
	// AlertRule represents an alert rule in LibreNMS.
	AlertRule struct {
		ID           int     `json:"id"`
		Builder      string  `json:"builder"`
		Devices      []int   `json:"devices"`
		Disabled     Bool    `json:"disabled"`
		Extra        string  `json:"extra"`
		Groups       []int   `json:"groups"`
		InvertMap    Bool    `json:"invert_map"`
		Locations    []int   `json:"locations"`
		Name         string  `json:"name"`
		Notes        *string `json:"notes"`
		ProcedureURL *string `json:"proc"`
		Query        string  `json:"query"`
		Rule         string  `json:"rule"`
		Severity     string  `json:"severity"`
	}

	// AlertRuleCreateRequest is the request structure for creating an alert rule.
	//
	// Groups and Locations can be empty, but Devices requires a -1 entry for 'all devices'.
	AlertRuleCreateRequest struct {
		Builder      string `json:"builder"`
		Count        int    `json:"count,omitempty"` // Max Alerts in the UI
		Delay        string `json:"delay,omitempty"`
		Devices      []int  `json:"devices"`
		Disabled     Bool   `json:"disabled,omitempty"`
		Groups       []int  `json:"groups"`
		Interval     string `json:"interval,omitempty"`
		Locations    []int  `json:"locations"`
		Mute         bool   `json:"mute,omitempty"`
		Name         string `json:"name"`
		Notes        string `json:"notes,omitempty"`
		ProcedureURL string `json:"proc,omitempty"`
		Query        string `json:"query,omitempty"`
		Rule         string `json:"rule,omitempty"`
		Severity     string `json:"severity"`
	}

	// AlertRuleUpdateRequest is the request structure for updating an alert rule.
	AlertRuleUpdateRequest struct {
		AlertRuleCreateRequest
		ID int `json:"rule_id"`
	}

	// AlertRuleResponse is the response structure for alert rules.
	AlertRuleResponse struct {
		BaseResponse
		Rules []AlertRule `json:"rules"`
	}
)

// CreateAlertRule creates a specific alert rule in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Alerts/#get_alert_rule
func (c *Client) CreateAlertRule(payload *AlertRuleCreateRequest) (*BaseResponse, error) {
	// as a convenience/hack, add a -1 to Devices if Devices is empty
	if len(payload.Devices) == 0 {
		payload.Devices = []int{-1}
	}

	req, err := c.newRequest(http.MethodPost, alertRuleEndpoint, payload, nil)
	if err != nil {
		return nil, err
	}
	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// DeleteAlertRule deletes a specific alert rule by its ID from the LibreNMS API.
func (c *Client) DeleteAlertRule(id int) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("%s/%d", alertRuleEndpoint, id), nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// GetAlertRule retrieves a specific alert rule by its ID from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Alerts/#get_alert_rule
func (c *Client) GetAlertRule(id int) (*AlertRuleResponse, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/%d", alertRuleEndpoint, id), nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(AlertRuleResponse)
	return resp, c.do(req, resp)
}

// GetAlertRules retrieves all alert rules from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Alerts/#list_alert_rules
func (c *Client) GetAlertRules() (*AlertRuleResponse, error) {
	req, err := c.newRequest(http.MethodGet, alertRuleEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(AlertRuleResponse)
	return resp, c.do(req, resp)
}

// UpdateAlertRule updates a specific alert rule in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Alerts/#update_alert_rule
func (c *Client) UpdateAlertRule(payload *AlertRuleUpdateRequest) (*BaseResponse, error) {
	if payload.ID < 1 {
		return nil, fmt.Errorf("rule ID is required for updating an alert rule")
	}

	// as a convenience/hack, add a -1 to Devices if Devices is empty
	if len(payload.Devices) == 0 {
		payload.Devices = []int{-1}
	}

	req, err := c.newRequest(http.MethodPut, alertRuleEndpoint, payload, nil)
	if err != nil {
		return nil, err
	}
	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}
