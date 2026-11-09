package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type VersionInfo struct {
	Version string `json:"version"`
	Env     string `json:"env"`
}

// GetPublicVersion returns the API version and environment name. Open
// (public, no bearer token) endpoint, GET /v3/public/version.
func (c *Client) GetPublicVersion(ctx context.Context) (*Response[VersionInfo], error) {
	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/v3/public/version",
		Public:   true,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[VersionInfo])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
