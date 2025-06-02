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
	//
	// Only set the fields you want to update. Trying to patch fields that have not changed will
	// result in an HTTP 500 error.
	ServiceUpdateRequest struct {
		Name        string `json:"service_name,omitempty"`
		Description string `json:"service_desc,omitempty"`
		IP          string `json:"service_ip,omitempty"`
		Ignore      Bool   `json:"service_ignore,omitempty"`
		Param       string `json:"service_param,omitempty"`
		Type        string `json:"service_type,omitempty"`
	}

	// serviceResponse is the internal response structure for services.
	//
	// The raw response is returned as a list of service lists, but it seems that all services are always returned
	// in the first list. This also causes the count to always reflect 1 in the response.
	// So ... we're going to collapse this into a 1-dimensional slice and update count for easier client handling.
	serviceResponse struct {
		BaseResponse
		Services [][]Service `json:"services"`
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
func (c *Client) CreateService(deviceIdentifier string, service *ServiceCreateRequest) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s/%s", serviceEndpoint, deviceIdentifier), service, nil)
	if err != nil {
		return nil, err
	}

	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}

// DeleteService deletes a service by its ID.
//
// Documentation: https://docs.librenms.org/API/Services/#delete_service_from_host
func (c *Client) DeleteService(serviceID int) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("%s/%d", serviceEndpoint, serviceID), nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(BaseResponse)
	err = c.do(req, resp)
	return resp, err
}

// GetService retrieves a service by ID from the LibreNMS API.
//
// Similar to GetDeviceGroup, this uses the same endpoint as GetServices, but it returns a
// modified payload with the single host. This is primarily a convenience function
// for the Terraform provider.
func (c *Client) GetService(serviceID int) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodGet, serviceEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	internalResp := new(serviceResponse)
	if err = c.do(req, internalResp); err != nil {
		return nil, err
	}

	resp := &ServiceResponse{
		BaseResponse: internalResp.BaseResponse,
		Services:     internalResp.getServices(),
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

	for _, service := range resp.Services {
		if service.ID == serviceID {
			singleServiceResp.Services = append(singleServiceResp.Services, service)
			singleServiceResp.Count = 1
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

	internalResp := new(serviceResponse)
	if err = c.do(req, internalResp); err != nil {
		return nil, err
	}

	services := internalResp.getServices()
	return &ServiceResponse{
		BaseResponse: BaseResponse{
			Status:  internalResp.Status,
			Message: internalResp.Message,
			Count:   len(services),
		},
		Services: services,
	}, err
}

// GetServicesForHost retrieves all services for a specific host by ID or name from the LibreNMS API.
//
// For whatever reason, there is no equivalent GetService endpoint.
// The /services/:id endpoint just returns the services for the host identifier. ¯\_(ツ)_/¯
// At least it's consistent with the device groups endpoint...
//
// Documentation: https://docs.librenms.org/API/Services/#get_service_for_host
func (c *Client) GetServicesForHost(deviceIdentifier string) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/%s", serviceEndpoint, deviceIdentifier), nil, nil)
	if err != nil {
		return nil, err
	}

	internalResp := new(serviceResponse)
	if err = c.do(req, internalResp); err != nil {
		return nil, err
	}

	services := internalResp.getServices()
	return &ServiceResponse{
		BaseResponse: BaseResponse{
			Status:  internalResp.Status,
			Message: internalResp.Message,
			Count:   len(services),
		},
		Services: services,
	}, err
}

// UpdateService updates a service for the specified service ID.
//
// Documentation: https://docs.librenms.org/API/Services/#edit_service_from_host
func (c *Client) UpdateService(serviceID int, service *ServiceUpdateRequest) (*ServiceResponse, error) {
	req, err := c.newRequest(http.MethodPatch, fmt.Sprintf("%s/%d", serviceEndpoint, serviceID), service, nil)
	if err != nil {
		return nil, err
	}

	resp := new(ServiceResponse)
	err = c.do(req, resp)
	return resp, err
}

// getServices flattens the slice of slices into a single slice
func (s *serviceResponse) getServices() []Service {
	flatServices := make([]Service, 0)
	for _, serviceList := range s.Services {
		flatServices = append(flatServices, serviceList...)
	}
	return flatServices
}
