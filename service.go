package librenms

import (
	"fmt"
	"net/http"
)

const (
	serviceEndpoint = "services"
)

type (
	// Service represents a service in LibreNMS.
	Service struct {
		ID          int    `json:"service_id"`
		Changed     int64  `json:"service_changed"`
		Description string `json:"service_desc"`
		DeviceID    int    `json:"device_id"`
		DS          string `json:"service_ds"` // I don't know what this is
		Ignore      Bool   `json:"service_ignore"`
		IP          string `json:"service_ip"`
		Message     string `json:"service_message"`
		Name        string `json:"service_name"`
		Param       string `json:"service_param"`
		Status      int    `json:"service_status"` // assuming this follows Nagios conventions, 0=ok, 1=warning, 2=critical, 3=unknown
		TemplateID  int    `json:"service_template_id"`
		Type        string `json:"service_type"`
	}

	// ServiceCreateRequest represents the request payload for creating a service.
	ServiceCreateRequest struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"desc,omitempty"`
		IP          string `json:"ip,omitempty"`
		Ignore      Bool   `json:"ignore,omitempty"`
		Param       string `json:"param,omitempty"`
		Type        string `json:"type"`
	}

	// ServiceUpdateRequest represents the request payload for updating a service.
	ServiceUpdateRequest struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"desc,omitempty"`
		IP          string `json:"ip,omitempty"`
		Ignore      Bool   `json:"ignore,omitempty"`
		Param       string `json:"param,omitempty"`
		Type        string `json:"type,omitempty"`
	}

	// ServiceResponse is the response structure for services.
	ServiceResponse struct {
		BaseResponse
		Services []Service `json:"services"`
	}
)

// CreateService creates a service for the specified device id or hostname.
//
// Documentation: https://docs.librenms.org/API/Services/#add_service_for_host
func (c *Client) CreateService(identifier string, service *ServiceCreateRequest) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s/%s", serviceEndpoint, identifier), service, nil)
	if err != nil {
		return nil, err
	}

	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}

// GetService retrieves a service by ID from the LibreNMS API.
//
// Similar to GetDeviceGroup, this uses the same endpoint as GetServices, but it returns a
// modified payload with the single host. This is primarily a convenience function
// for the Terraform provider.
func (c *Client) GetService(id int) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodGet, serviceEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(ServiceResponse)
	err = c.do(req, resp)
	if err != nil {
		return resp, err
	}

	if len(resp.Services) == 0 {
		return resp, nil
	}

	// look for a matching service by ID
	singleServiceResp := &ServiceResponse{
		Services: make([]Service, 0),
	}
	singleServiceResp.Message = resp.Message
	singleServiceResp.Status = resp.Status
	singleServiceResp.Count = 1

	for _, service := range resp.Services {
		if service.ID == id {
			singleServiceResp.Services = append(singleServiceResp.Services, service)
			break
		}
	}

	return singleServiceResp, err
}

// GetServices retrieves all services from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Services/#list_services
func (c *Client) GetServices() (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodGet, serviceEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}

// GetServicesForHost retrieves all services for a specific host by ID or name from the LibreNMS API.
//
// For whatever reason, there is no equivalent GetService endpoint.
// The /services/:id endpoint just returns the services for the host identifier. ¯\_(ツ)_/¯
// At least it's consistent with the device groups endpoint...
//
// Documentation: https://docs.librenms.org/API/Services/#get_service_for_host
func (c *Client) GetServicesForHost(identifier string) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodGet, serviceEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}

// UpdateService updates a service for the specified device id or hostname.
//
// Documentation: https://docs.librenms.org/API/Services/#add_service_for_host
func (c *Client) UpdateService(identifier string, service *ServiceUpdateRequest) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s/%s", serviceEndpoint, identifier), service, nil)
	if err != nil {
		return nil, err
	}

	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}
