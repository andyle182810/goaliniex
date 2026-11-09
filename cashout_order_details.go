package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetCashoutOrder retrieves one cashout order by the orderId the caller
// assigned, GET /v3/cashout/orders/{orderId}.
func (c *Client) GetCashoutOrder(
	ctx context.Context,
	userEmail string,
	orderID string,
) (*Response[CashoutOrder], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/cashout/orders/" + url.PathEscape(orderID),
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[CashoutOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
