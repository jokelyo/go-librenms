// Package librenms provides a client for interacting with the LibreNMS API.
//
// The package supports CRUD operations for the following resources, which
// are used within the Terraform provider at https://github.com/jokelyo/terraform-provider-librenms:
//   - Alert Rules
//   - Devices
//   - Device Groups
//   - Services
//
// LibreNMS API Documentation: https://docs.librenms.org/API/
package librenms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
)

const (
	apiVersion = "v0"
	authHeader = "X-Auth-Token"
)

type (
	// Bool represents a boolean value, used for JSON marshaling. The API
	// returns some fields as 0/1 instead of true/false, so we use this custom type.
	Bool bool

	// Client is the main structure for the LibreNMS client.
	Client struct {
		baseURL *url.URL
		client  *http.Client
		token   string
	}

	// Option is a function that configures the Client.
	Option func(*Client)

	// BaseResponse is the base structure for API responses.
	BaseResponse struct {
		// Status indicates the success or failure of the API call.
		Status string `json:"status"`
		// Message contains additional information about the API call.
		Message string `json:"message"`
		Count   int    `json:"count"`
	}
)

// WithHTTPClient sets the HTTP client for the LibreNMS client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// New creates a new LibreNMS client with the given base URL and options.
// The base URL should be in the format 'http[s]://<host>[:port]/'.
func New(baseURL, token string, opts ...Option) (*Client, error) {
	c := &Client{
		token:  token,
		client: cleanhttp.DefaultPooledClient(),
	}

	// Append a trailing slash to the base URL if it doesn't have one.
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// Parse the base URL
	var err error
	c.baseURL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if c.baseURL.Path != "/" {
		return nil, errors.New("invalid base URL format, expected: 'http[s]://<host>[:port]/'")
	}

	// Append the API version to the base URL path.
	c.baseURL, err = c.baseURL.Parse(fmt.Sprintf("api/%s/", apiVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to parse API version in base URL: %w", err)
	}

	// Process options
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// newRequest creates a new HTTP request with the given method and path.
// A relative URI should be provided and should not have a leading slash.
func (c *Client) newRequest(method, uri string, body any, query *url.Values) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	// Parse the URI and construct the full URL
	fullURL, err := c.baseURL.Parse(uri)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(context.Background(), method, fullURL.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set necessary headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set(authHeader, c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add query parameters if provided
	if query != nil && len(*query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

// rawDo sends an HTTP request and returns the raw response body. We should normally
// use do() which JSON-decodes and closes the response body, but if there is a non-JSON
// endpoint or other reason to not decode, this can be used.
func (c *Client) rawDo(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, checkResponse(resp)
}

// do sends an HTTP request and decodes the JSON response into the provided response object.
func (c *Client) do(req *http.Request, respObj any) error {
	if respObj == nil {
		return errors.New("response object cannot be nil")
	}

	resp, err := c.rawDo(req)
	defer closeBody(resp.Body)
	if err != nil {
		return err
	}

	switch v := respObj.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if errors.Is(decErr, io.EOF) {
			decErr = nil // No content to decode, treat as success
		}
		if decErr != nil {
			err = fmt.Errorf("failure decoding response: %w", decErr)
		}
	}
	return err
}

// checkResponse checks the HTTP response for errors.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorResponse := &ErrorResponse{
			Response: resp,
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errorResponse.Message = fmt.Sprintf("failed to read response body: %v", err)
			return errorResponse
		}

		if len(body) > 0 {
			if err = json.Unmarshal(body, errorResponse); err != nil {
				errorResponse.Message = string(body)
			}
		}

		return errorResponse
	}
	return nil
}

func closeBody(body io.ReadCloser) {
	_ = body.Close()
}

// parseParams is a helper function that parses the provided value into URL query parameters.
func parseParams(v any) (*url.Values, error) {
	if v == nil {
		return new(url.Values), nil
	}

	p, err := query.Values(v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query parameters: %w", err)
	}
	return &p, nil
}

// MarshalJSON implements the JSON marshaling for the Bool type.
func (b *Bool) MarshalJSON() ([]byte, error) {
	if *b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// UnmarshalJSON implements the JSON unmarshalling for the Bool type.
func (b *Bool) UnmarshalJSON(data []byte) error {
	// attempt to unmarshal as a boolean first
	var valueBool bool
	if err := json.Unmarshal(data, &valueBool); err == nil {
		*b = Bool(valueBool)
		return nil
	}

	// if that fails, try to unmarshal as an integer (0 or 1)
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("failed to unmarshal Bool: %w", err)
	}
	*b = value != 0
	return nil
}
