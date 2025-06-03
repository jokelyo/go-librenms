package librenms_test

import (
	"fmt"
	"log"
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
)

// This init function will register handlers for location-related API endpoints.
func init() {
	mockCreateLocationResponse := loadMockResponse("create_location_200.json")
	mockUpdateLocationResponse := loadMockResponse("update_location_200.json")
	mockDeleteLocationResponse := loadMockResponse("delete_location_200.json")
	mockGetLocationsResponse := loadMockResponse("get_locations_200.json")
	mockGetLocationResponse := loadMockResponse("get_location_200.json")

	// Handle GET for locations resource
	mux.HandleFunc(testEndpointLocations, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetLocationsResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", r.Method, testEndpointLocations), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Handle GET for a specific location
	mux.HandleFunc(testEndpointLocation, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write(mockGetLocationResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", r.Method, testEndpointLocation), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Handle POST for creating locations
	mux.HandleFunc("/api/v0/locations", func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			_, err = w.Write(mockCreateLocationResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for /api/v0/locations.", r.Method), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Handle PATCH and DELETE for specific location
	mux.HandleFunc(testEndpointLocationPatch, func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPatch:
			_, err = w.Write(mockUpdateLocationResponse)
		case http.MethodDelete:
			_, err = w.Write(mockDeleteLocationResponse)
		default:
			http.Error(w, fmt.Sprintf("Method %s not implemented for %s.", r.Method, testEndpointLocationPatch), http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func TestClient_GetLocation(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	locationResp, err := testAPIClient.GetLocation(testLocationID)

	r.NoError(err, "GetLocation returned an error")
	r.NotNil(locationResp, "GetLocation response is nil")

	r.Equal("ok", locationResp.Status, "Expected status 'ok'")

	r.Equal(2, locationResp.Location.ID, "Expected location ID 2")
	r.Equal(librenms.Bool(false), locationResp.Location.FixedCoordinates, "Expected FixedCoordinates to be false")
}

func TestClient_GetLocations(t *testing.T) {
	r := require.New(t)

	r.NotNil(testAPIClient, "Global testAPIClient should be initialized")

	locationResp, err := testAPIClient.GetLocations()

	r.NoError(err, "GetLocations returned an error")
	r.NotNil(locationResp, "GetLocations response is nil")

	r.Equal("ok", locationResp.Status, "Expected status 'ok'")
	r.Equal(2, locationResp.Count, "Expected count 2")
	r.Len(locationResp.Locations, 2, "Expected 2 locations")

	location := locationResp.Locations[0]
	r.Equal(1, location.ID, "Expected Location ID 1")
	r.Equal("Sitting on the Dock of the Bay", location.Name, "Expected Location name 'Sitting on the Dock of the Bay'")
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
