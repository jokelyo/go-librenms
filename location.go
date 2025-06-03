package librenms

import (
	"fmt"
	"net/http"
)

type (
	// Location represents a location in LibreNMS.
	Location struct {
		ID               int     `json:"id"`
		FixedCoordinates Bool    `json:"fixed_coordinates"`
		Latitude         float64 `json:"lat"`
		Longitude        float64 `json:"lng"`
		Name             string  `json:"location"`
		Timestamp        string  `json:"timestamp"`
	}

	// LocationCreateRequest represents the request payload for creating a location.
	LocationCreateRequest struct {
		Name             string  `json:"location"`
		FixedCoordinates Bool    `json:"fixed_coordinates"`
		Latitude         float64 `json:"lat"`
		Longitude        float64 `json:"lng"`
	}

	// LocationUpdateRequest represents the request payload for updating a location.
	//
	// Only set the field(s) you want to update. Trying to patch fields that have not changed will
	// result in an HTTP 500 error.
	LocationUpdateRequest struct {
		Name             *string
		FixedCoordinates *bool
		Latitude         *float64
		Longitude        *float64
	}

	// LocationResponse represents a response containing a single location from the LibreNMS API.
	LocationResponse struct {
		Status   string   `json:"status"`
		Location Location `json:"get_location"`
	}

	// LocationResponse represents a response containing a list of locations from the LibreNMS API.
	LocationsResponse struct {
		BaseResponse
		Locations []Location `json:"locations"`
	}
)

// CreateLocation creates a new location in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Locations/#add_location
func (c *Client) CreateLocation(location *LocationCreateRequest) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodPost, "locations", location, nil)
	if err != nil {
		return nil, err
	}

	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// DeleteLocation deletes a location by its ID in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Locations/#delete_location
func (c *Client) DeleteLocation(locationID int) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("locations/%d", locationID), nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// GetLocation retrieves a location by its ID from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Locations/#get_location
func (c *Client) GetLocation(locationID int) (*LocationResponse, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("location/%d", locationID), nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(LocationResponse)
	return resp, c.do(req, resp)
}

// GetLocations retrieves a list of locations from the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Locations/#list_locations
func (c *Client) GetLocations() (*LocationsResponse, error) {
	req, err := c.newRequest(http.MethodGet, "resources/locations", nil, nil)
	if err != nil {
		return nil, err
	}

	resp := new(LocationsResponse)
	return resp, c.do(req, resp)
}

// UpdateLocation updates a location by its ID in the LibreNMS API.
//
// Documentation: https://docs.librenms.org/API/Locations/#edit_location
func (c *Client) UpdateLocation(locationID int, location *LocationUpdateRequest) (*BaseResponse, error) {
	req, err := c.newRequest(http.MethodPatch, fmt.Sprintf("locations/%d", locationID), location.payload(), nil)
	if err != nil {
		return nil, err
	}

	resp := new(BaseResponse)
	return resp, c.do(req, resp)
}

// NewLocationUpdateRequest creates a new, empty LocationUpdateRequest.
func NewLocationUpdateRequest() *LocationUpdateRequest {
	return &LocationUpdateRequest{}
}

// SetFixedCoordinates sets whether the location has fixed coordinates in the LocationUpdateRequest.
func (r *LocationUpdateRequest) SetFixedCoordinates(fixed bool) *LocationUpdateRequest {
	r.FixedCoordinates = &fixed
	return r
}

// SetLatitude sets the latitude of the location in the LocationUpdateRequest.
func (r *LocationUpdateRequest) SetLatitude(lat float64) *LocationUpdateRequest {
	r.Latitude = &lat
	return r
}

// SetLongitude sets the longitude of the location in the LocationUpdateRequest.
func (r *LocationUpdateRequest) SetLongitude(lng float64) *LocationUpdateRequest {
	r.Longitude = &lng
	return r
}

// SetName sets the name of the location in the LocationUpdateRequest.
func (r *LocationUpdateRequest) SetName(name string) *LocationUpdateRequest {
	r.Name = &name
	return r
}

// payload converts the LocationUpdateRequest to a map for the API request.
func (r *LocationUpdateRequest) payload() map[string]interface{} {
	payload := make(map[string]interface{})
	if r.Name != nil {
		payload["location"] = *r.Name
	}
	if r.FixedCoordinates != nil {
		payload["fixed_coordinates"] = *r.FixedCoordinates
	}
	if r.Latitude != nil {
		payload["lat"] = *r.Latitude
	}
	if r.Longitude != nil {
		payload["lng"] = *r.Longitude
	}
	return payload
}
