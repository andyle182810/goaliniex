package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetScanToPayOrder retrieves one scan-to-pay order by the orderId the
// caller assigned, GET /v3/scan-to-pay/orders/{orderId}.
func (c *Client) GetScanToPayOrder(
	ctx context.Context,
	userEmail string,
	orderID string,
) (*Response[ScanToPayOrder], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/scan-to-pay/orders/" + url.PathEscape(orderID),
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[ScanToPayOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
