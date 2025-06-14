package librenms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// deviceGroupEndpoint is the API endpoint for devices.
	deviceGroupEndpoint = "devicegroups"
)

type (
	// DeviceGroup represents a device group in LibreNMS.
	DeviceGroup struct {
		ID          int                      `json:"id"`
		Name        string                   `json:"name"`
		Description *string                  `json:"desc"`
		Pattern     *string                  `json:"pattern"`
		Rules       DeviceGroupRuleContainer `json:"rules"`
		Type        string                   `json:"type"`
	}

	// DeviceGroupRuleContainer represents the top-level container for device group rules.
	DeviceGroupRuleContainer struct {
		Condition string            `json:"condition"`
		Joins     [][]string        `json:"joins"`
		Rules     []DeviceGroupRule `json:"rules"`
		Valid     bool              `json:"valid""`
	}

	// DeviceGroupRule represents a rule within a device group. This is a recursive structure.
	// It can contain nested rules, allowing for complex conditions.
	//
	// A terminal section defines id, field, type, input, operator, and value.
	// A non-terminal section defines condition and a list of rules.
	DeviceGroupRule struct {
		ID        string            `json:"id,omitempty"`
		Condition string            `json:"condition,omitempty"`
		Field     string            `json:"field,omitempty"`
		Input     string            `json:"input,omitempty"`
		Operator  string            `json:"operator,omitempty"`
		Rules     []DeviceGroupRule `json:"rules,omitempty"`
		Type      string            `json:"type,omitempty"`
		Value     string            `json:"value,omitempty"`
	}

	// DeviceGroupCreateRequest represents the request payload for creating a device group.
	//
	// The rules should be a serialized JSON string that matches the DeviceGroupRuleContainer
	// structure. Define your rules using the DeviceGroupRuleContainer struct and then
	// serialize it using its JSON() method.
	DeviceGroupCreateRequest struct {
		Name        string  `json:"name"`
		Description *string `json:"desc,omitempty"`
		Devices     []int   `json:"devices,omitempty"`
		Rules       *string `json:"rules,omitempty"`
		Type        string  `json:"type"`
	}

	// DeviceGroupUpdateRequest represents the request payload for updating a device group.
	//
	// The rules should be a serialized JSON string that matches the DeviceGroupRuleContainer
	// structure. Define your rules using the DeviceGroupRuleContainer struct and then
	// serialize it using its JSON() method.
	DeviceGroupUpdateRequest struct {
		Name        string  `json:"name,omitempty"`
		Description *string `json:"desc,omitempty"`
		Devices     []int   `json:"devices,omitempty"`
		Rules       *string `json:"rules,omitempty"`
		Type        string  `json:"type,omitempty"`
	}

	// DeviceGroupResponse represents a response containing a list of device groups from the LibreNMS API.
	DeviceGroupResponse struct {
		BaseResponse
		Groups []DeviceGroup `json:"groups"`
	}

	// DeviceGroupMember represents a member of a device group.
	DeviceGroupMember struct {
		ID int `json:"device_id"`
	}

	// DeviceGroupMembersResponse represents a response containing the members of a device group.
	DeviceGroupMembersResponse struct {
		BaseResponse
		Devices []DeviceGroupMember `json:"devices"`
	}

	// DeviceGroupCreateResponse represents a creation response.
	DeviceGroupCreateResponse struct {
		BaseResponse
		ID int `json:"id"`
	}
)

// CreateDeviceGroup creates a device group in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/DeviceGroups/#add_devicegroup
func (c *Client) CreateDeviceGroup(group *DeviceGroupCreateRequest) (*DeviceGroupCreateResponse, error) {
	req, err := c.newRequest(http.MethodPost, deviceGroupEndpoint, group, nil)
	if err != nil {
		return nil, err
	}

	resp := new(DeviceGroupCreateResponse)
	return resp, c.do(req, resp)
}

// DeleteDeviceGroup deletes a group by its ID or hostname from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/DeviceGroups/#delete_devicegroup
func (c *Client) DeleteDeviceGroup(identifier string) (*BaseResponse, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/%s", deviceGroupEndpoint, identifier))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	req, err := c.newRequest(http.MethodDelete, uri.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// GetDeviceGroup uses the same endpoint as GetDeviceGroups, but it returns a
// modified payload with the single host (if a match is found).
// This is primarily a convenience function for the Terraform provider.
func (c *Client) GetDeviceGroup(identifier string) (*DeviceGroupResponse, error) {
	req, err := c.newRequest(http.MethodGet, deviceGroupEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(DeviceGroupResponse)
	if err = c.do(req, resp); err != nil {
		return resp, err
	}

	if len(resp.Groups) == 0 {
		return resp, nil
	}

	singleGroupResp := &DeviceGroupResponse{
		Groups: make([]DeviceGroup, 0),
	}
	singleGroupResp.Message = resp.Message
	singleGroupResp.Status = resp.Status

	for _, group := range resp.Groups {
		if group.Name == identifier || strconv.Itoa(group.ID) == identifier {
			singleGroupResp.Groups = append(singleGroupResp.Groups, group)
			singleGroupResp.Count = 1
			break
		}
	}

	return singleGroupResp, err
}

// GetDeviceGroups retrieves a list of device groups from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/DeviceGroups/#get_devicegroups
func (c *Client) GetDeviceGroups() (*DeviceGroupResponse, error) {
	req, err := c.newRequest(http.MethodGet, deviceGroupEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(DeviceGroupResponse)
	return resp, c.do(req, resp)
}

// GetDeviceGroupMembers retrieves a list of device group members from the LibreNMS API.
// The identifier can be either the group ID or the group name.
//
// Documentation: https://docs.librenms.org/API/DeviceGroups/#get_devices_by_group
func (c *Client) GetDeviceGroupMembers(identifier string) (*DeviceGroupMembersResponse, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/%s", deviceGroupEndpoint, identifier), nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(DeviceGroupMembersResponse)
	return resp, c.do(req, resp)
}

// UpdateDeviceGroup updates an existing device group in the LibreNMS API.
//
// The documentation states it uses name rather than ID to reference the group, but both seem to work (as of v25.5).
// Documentation: https://docs.librenms.org/API/DeviceGroups/#update_devicegroup
func (c *Client) UpdateDeviceGroup(identifier string, payload *DeviceGroupUpdateRequest) (*BaseResponse, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/%s", deviceGroupEndpoint, identifier))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	req, err := c.newRequest(http.MethodPatch, uri.String(), payload, nil)
	if err != nil {
		return nil, err
	}

	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// JSON is a helper function that serializes the DeviceGroupRuleContainer to JSON format.
func (g *DeviceGroupRuleContainer) JSON() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MustJSON is a helper function that serializes the DeviceGroupRuleContainer to JSON format.
// It returns an empty string if the marshalling fails.
func (g *DeviceGroupRuleContainer) MustJSON() string {
	data, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(data)
}
