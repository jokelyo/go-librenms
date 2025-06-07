package librenms_test

import (
	"net/http"
	"testing"

	"github.com/jokelyo/go-librenms"

	"github.com/stretchr/testify/require"
)

const (
	testLocationID            = 2
	testEndpointLocations     = "/api/v0/resources/locations"
	testEndpointLocation      = "/api/v0/location/2"
	testEndpointLocationPatch = "/api/v0/locations/2"
	testEndpointLocationsPost = "/api/v0/locations"
)

// This init function will register handlers for location-related API endpoints.
func init() {
	handleEndpoint(testEndpointLocations, mockResponses{
		http.MethodGet: loadMockResponse("get_locations_200.json"),
	})

	handleEndpoint(testEndpointLocation, mockResponses{
		http.MethodGet: loadMockResponse("get_location_200.json"),
	})

	handleEndpoint(testEndpointLocationsPost, mockResponses{
		http.MethodPost: loadMockResponse("create_location_200.json"),
	})

	handleEndpoint(testEndpointLocationPatch, mockResponses{
		http.MethodPatch:  loadMockResponse("update_location_200.json"),
		http.MethodDelete: loadMockResponse("delete_location_200.json"),
	})
}

func TestClient_GetLocation(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	locationResp, err := testAPIClient.GetLocation(testLocationID)

	r.NoError(err, "GetLocation returned an error")
	r.NotNil(locationResp, "GetLocation response is nil")

	r.Equal("ok", locationResp.Status, "Expected status 'ok'")

	r.Equal(1, locationResp.Location.ID, "Expected location ID 1")
	r.Equal("test location", locationResp.Location.Name, "Expected Location name 'test location'")
	r.Equal(librenms.Bool(true), locationResp.Location.FixedCoordinates, "Expected FixedCoordinates to be true")
	r.Equal(librenms.Float64(37.4220648), locationResp.Location.Longitude, "Expected Longitude to be 37.4220648")
}

func TestClient_GetLocations(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	locationResp, err := testAPIClient.GetLocations()

	r.NoError(err, "GetLocations returned an error")
	r.NotNil(locationResp, "GetLocations response is nil")

	r.Equal("ok", locationResp.Status, "Expected status 'ok'")
	r.Equal(5, locationResp.Count, "Expected count 5")
	r.Len(locationResp.Locations, 5, "Expected 5 locations")

	location := locationResp.Locations[0]
	r.Equal(1, location.ID, "Expected Location ID 1")
	r.Equal("Sitting on the Dock of the Bay", location.Name, "Expected Location name 'Sitting on the Dock of the Bay'")
	r.Equal(librenms.Bool(false), location.FixedCoordinates, "Expected FixedCoordinates to be false")

	location = locationResp.Locations[4]
	r.Equal(5, location.ID, "Expected Location ID 5")
	r.Equal("test location", location.Name, "Expected Location name 'test location'")
	r.Equal(librenms.Bool(true), location.FixedCoordinates, "Expected FixedCoordinates to be true")
	r.Equal(librenms.Float64(37.42206480), location.Longitude, "Expected Longitude to be 37.4220648")
}

func TestClient_CreateLocation(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	newLocationRequest := librenms.LocationCreateRequest{
		Name:             "Test Location",
		FixedCoordinates: false,
		Latitude:         37.7749,
		Longitude:        -122.4194,
	}

	createResp, err := testAPIClient.CreateLocation(&newLocationRequest)

	r.NoError(err, "CreateLocation returned an error")
	r.NotNil(createResp, "CreateLocation response is nil")

	r.Equal("ok", createResp.Status, "Expected status 'ok'")
}

func TestClient_DeleteLocation(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	resp, err := testAPIClient.DeleteLocation(testLocationID)

	r.NoError(err, "DeleteLocation returned an error")
	r.NotNil(resp, "DeleteLocation response is nil")

	r.Equal("ok", resp.Status, "Expected status 'ok'")
}

func TestClient_UpdateLocation(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	updateResp, err := testAPIClient.UpdateLocation(
		testLocationID,
		librenms.NewLocationUpdateRequest().SetName("Updated Test Location"),
	)

	r.NoError(err, "UpdateLocation returned an error")
	r.NotNil(updateResp, "UpdateLocation response is nil")

	r.Equal("ok", updateResp.Status, "Expected status 'ok'")
}
