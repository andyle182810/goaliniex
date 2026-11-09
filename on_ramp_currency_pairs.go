package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetOnRampCurrencyPairsRequest filters are all optional and matched
// case-insensitively by the server.
type GetOnRampCurrencyPairsRequest struct {
	FiatCurrency   string `json:"fiatCurrency,omitempty"`
	CryptoCurrency string `json:"cryptoCurrency,omitempty"`
	Network        string `json:"network,omitempty"`
}

// GetOnRampCurrencyPairs lists the fiat/crypto/network combinations that
// can be bought via on-ramp, GET /v3/on-ramp/currency-pairs.
func (c *Client) GetOnRampCurrencyPairs(
	ctx context.Context,
	userEmail string,
	req *GetOnRampCurrencyPairsRequest,
) (*Response[[]CurrencyPair], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/on-ramp/currency-pairs",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]CurrencyPair])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
