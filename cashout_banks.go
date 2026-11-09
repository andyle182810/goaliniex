package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetCashoutBanks lists the banks payable via cashout for fiatCurrency,
// GET /v3/cashout/banks.
func (c *Client) GetCashoutBanks(
	ctx context.Context,
	userEmail string,
	fiatCurrency FiatCurrency,
) (*Response[[]Bank], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/cashout/banks",
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
