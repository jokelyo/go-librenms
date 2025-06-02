package librenms

import (
	"fmt"
	"net/http"
)

const (
	// deviceEndpoint is the API endpoint for devices.
	deviceEndpoint = "devices"
)

type (
	// Device represents a device in LibreNMS.
	//
	// Pointers are used for fields that may be null.
	// A custom type Bool is used to represent booleans that may be defined as 0/1 by the API.
	Device struct {
		DeviceID int `json:"device_id"`

		AgentUptime             int      `json:"agent_uptime"`
		AuthAlgorithm           *string  `json:"authalgo"`
		AuthLevel               *string  `json:"authlevel"`
		AuthName                *string  `json:"authname"`
		AuthPass                *string  `json:"authpass"`
		BGPLocalAS              *int     `json:"bgpLocalAs"`
		Community               *string  `json:"community"`
		CryptoAlgorithm         *string  `json:"cryptoalgo"`
		CryptoPass              *string  `json:"cryptopass"`
		DisableNotify           Bool     `json:"disable_notify"`
		Disabled                Bool     `json:"disabled"`
		Display                 *string  `json:"display"`
		Features                *string  `json:"features"`
		Hardware                string   `json:"hardware"`
		Hostname                string   `json:"hostname"`
		Icon                    string   `json:"icon"`
		Ignore                  Bool     `json:"ignore"`
		IgnoreStatus            Bool     `json:"ignore_status"`
		Inserted                string   `json:"inserted"`
		IP                      string   `json:"ip"`
		LastDiscovered          *string  `json:"last_discovered"`
		LastDiscoveredTimeTaken float64  `json:"last_discovered_timetaken"`
		LastPing                *string  `json:"last_ping"`
		LastPingTimeTaken       float64  `json:"last_ping_timetaken"`
		LastPollAttempted       *string  `json:"last_poll_attempted"`
		LastPolled              *string  `json:"last_pulled"`
		LastPolledTimeTaken     float64  `json:"last_polled_timetaken"`
		Latitude                *float64 `json:"lat"`
		Longitude               *float64 `json:"lng"`
		Location                *string  `json:"location"`
		LocationID              *int     `json:"location_id"`
		MaxDepth                *int     `json:"max_depth"`
		Notes                   *string  `json:"notes"`
		OS                      string   `json:"os"`
		OverrideSysLocation     Bool     `json:"override_sysLocation"`
		OverwriteIP             string   `json:"overwrite_ip"`
		PollerGroup             int      `json:"poller_group"`
		Port                    int      `json:"port"`
		PortAssociationMode     int      `json:"port_association_mode"`
		Purpose                 *string  `json:"purpose"`
		Retries                 *int     `json:"retries"`
		Serial                  *string  `json:"serial"`
		SNMPDisable             Bool     `json:"snmp_disable"`
		SNMPVersion             string   `json:"snmpver"`
		Status                  Bool     `json:"status"` // /devices returns 0/1, and /devices/:id returns true/false
		StatusReason            string   `json:"status_reason"`
		SysContact              *string  `json:"sysContact"`
		SysDescr                *string  `json:"sysDescr"`
		SysName                 string   `json:"sysName"`
		SysObjectID             *string  `json:"sysObjectID"`
		Timeout                 *int     `json:"timeout"`
		Transport               string   `json:"transport"`
		Type                    string   `json:"type"`
		Uptime                  *int64   `json:"uptime"`
		Version                 *string  `json:"version"`
	}

	// DeviceCreateRequest represents the request body for creating a new device in LibreNMS.
	DeviceCreateRequest struct {
		Hostname            string `json:"hostname"`
		Display             string `json:"display,omitempty"`
		ForceAdd            bool   `json:"force_add,omitempty"`
		Hardware            string `json:"hardware,omitempty"`
		Location            string `json:"location,omitempty"`
		LocationID          int    `json:"location_id,omitempty"`
		OS                  string `json:"os,omitempty"`
		OverrideSysLocation bool   `json:"override_sysLocation,omitempty"`
		PingFallback        bool   `json:"ping_fallback,omitempty"`
		PollerGroup         int    `json:"poller_group,omitempty"`
		Port                int    `json:"port,omitempty"`
		PortAssocMode       int    `json:"port_association_mode,omitempty"` // ifIndex(1), ifName(2), ifDescr(3), ifAlias(4)
		SNMPAuthAlgo        string `json:"authalgo,omitempty"`              // MD5, SHA, SHA-224, SHA-256, SHA384, SHA-512
		SNMPAuthLevel       string `json:"authlevel,omitempty"`             // noAuthNoPriv, authNoPriv, authPriv
		SNMPAuthName        string `json:"authname,omitempty"`
		SNMPAuthPass        string `json:"authpass,omitempty"`
		SNMPCrytoAlgo       string `json:"cryptoalgo,omitempty"` // DES, AES, AES-192, AES-256, AES-256-C
		SNMPCryptoPass      string `json:"cryptopass,omitempty"`
		SNMPCommunity       string `json:"community,omitempty"`
		SNMPDisable         bool   `json:"snmp_disable,omitempty"`
		SNMPVersion         string `json:"snmpver,omitempty"` // v1, v2c, v3
		SysName             string `json:"sysName,omitempty"`
		Transport           string `json:"transport,omitempty"`
	}

	// DeviceUpdateRequest represents the request body for updating a device in LibreNMS.
	//
	// The `Field` slice contains the names of the field(s) to update,
	// and `Data` contains the corresponding values. Only specify the fields you want to update.
	DeviceUpdateRequest struct {
		Field []string `json:"field"`
		Data  []any    `json:"data"`
	}

	// DeviceResponse represents a response containing a list of devices from the LibreNMS API.
	DeviceResponse struct {
		BaseResponse
		Devices []Device `json:"devices"`
	}

	DevicesQuery struct {
		DeviceID   int    `url:"device_id,omitempty"`
		Display    string `url:"display,omitempty"`
		Hostname   string `url:"hostname,omitempty"`
		IPv4       string `url:"ipv4,omitempty"`
		IPv6       string `url:"ipv6,omitempty"`
		Location   string `url:"location,omitempty"`
		LocationID int    `url:"location_id,omitempty"`
		MACAddress string `url:"mac,omitempty"`
		Order      string `url:"order,omitempty"`
		OS         string `url:"os,omitempty"`
		SysName    string `url:"sysName,omitempty"`
		Type       string `url:"type,omitempty"`
	}
)

// CreateDevice creates a device by hostname/IP.
//
// Documentation: https://docs.librenms.org/API/Devices/#add_device
func (c *Client) CreateDevice(payload *DeviceCreateRequest) (*DeviceResponse, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s/", deviceEndpoint), payload, nil)
	if err != nil {
		return nil, err
	}
	deviceResp := new(DeviceResponse)
	return deviceResp, c.do(req, deviceResp)
}

// DeleteDevice deletes a device by its ID or hostname from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Devices/#del_device
func (c *Client) DeleteDevice(identifier string) (*DeviceResponse, error) {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("%s/%s", deviceEndpoint, identifier), nil, nil)
	if err != nil {
		return nil, err
	}
	deviceResp := new(DeviceResponse)
	return deviceResp, c.do(req, deviceResp)
}

// GetDevice retrieves a device by its ID or hostname from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Devices/#get_device
func (c *Client) GetDevice(identifier string) (*DeviceResponse, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/%s", deviceEndpoint, identifier), nil, nil)
	if err != nil {
		return nil, err
	}
	deviceResp := new(DeviceResponse)
	return deviceResp, c.do(req, deviceResp)
}

// GetDevices retrieves a list of devices from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Devices/#list_devices
func (c *Client) GetDevices(query *DevicesQuery) (*DeviceResponse, error) {
	params, err := parseParams(query)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodGet, deviceEndpoint, nil, params)
	if err != nil {
		return nil, err
	}

	deviceResp := new(DeviceResponse)
	return deviceResp, c.do(req, deviceResp)
}

// UpdateDevice updates a device by hostname/IP.
//
// Documentation: https://docs.librenms.org/API/Devices/#update_device_field
func (c *Client) UpdateDevice(identifier string, payload *DeviceUpdateRequest) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodPatch, fmt.Sprintf("%s/%s", deviceEndpoint, identifier), payload, nil)
	if err != nil {
		return nil, err
	}
	patchResp := new(BaseResponse)
	return patchResp, c.do(req, patchResp)
}
