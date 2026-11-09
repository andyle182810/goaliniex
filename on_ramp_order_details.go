package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetOnRampOrder retrieves one on-ramp order by the orderId the caller
// assigned, GET /v3/on-ramp/orders/{orderId}.
func (c *Client) GetOnRampOrder(
	ctx context.Context,
	userEmail string,
	orderID string,
) (*Response[OnRampOrder], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/on-ramp/orders/" + url.PathEscape(orderID),
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[OnRampOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
