package librenms

import (
	"fmt"
	"net/http"
)

type (
	// Device represents a device in LibreNMS. Pointers are used for fields that may be null.
	// A custom type Bool is used to represent booleans that are defined as 0/1 by the API.
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
		Status                  bool     `json:"status"`
		StatusReason            string   `json:"status_reason"`
		SysContact              *string  `json:"sysContact"`
		SysDescr                *string  `json:"sysDescr"`
		SysName                 string   `json:"sysName"`
		SysObjectID             *string  `json:"sysObjectID"`
		Timeout                 *int     `json:"timeout"`
		Transport               string   `json:"transport"`
		Type                    string   `json:"type"`
		Uptime                  *string  `json:"uptime"`
		Version                 *string  `json:"version"`
	}

	// DeviceResponse represents a response containing a list of devices from the LibreNMS API.
	DeviceResponse struct {
		BaseResponse
		Devices []Device `json:"devices"`
	}
)

// GetDevice retrieves a device by its ID or hostname from the LibreNMS API.
func (c *Client) GetDevice(identifier string) (*DeviceResponse, error) {
	req, err := c.newRequest(http.MethodGet, "devices/"+identifier, nil, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("GetDevice: ", req.URL.String())
	fmt.Println("GetDevice: ", req.Method)
	for k, v := range req.Header {
		fmt.Printf("GetDevice: %s: %v\n", k, v)
		fmt.Println()
	}
	deviceResp := new(DeviceResponse)
	err = c.do(req, deviceResp)
	return deviceResp, err
}
