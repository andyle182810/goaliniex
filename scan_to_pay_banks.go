package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetScanToPayBanks lists the banks payable via scan-to-pay for fiatCurrency,
// GET /v3/scan-to-pay/banks.
func (c *Client) GetScanToPayBanks(
	ctx context.Context,
	userEmail string,
	fiatCurrency FiatCurrency,
) (*Response[[]Bank], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/scan-to-pay/banks",
		Params:    map[string]any{"fiatCurrency": fiatCurrency},
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]Bank])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
