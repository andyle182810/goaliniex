package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// CurrencyPair is a fiat/crypto/network combination available for off-ramp
// or on-ramp, with its amount limits (in crypto). Shared by both products.
type CurrencyPair struct {
	FiatCurrency   string  `json:"fiatCurrency"`
	CryptoCurrency string  `json:"cryptoCurrency"`
	Network        string  `json:"network"`
	NetworkName    string  `json:"networkName"`
	Fee            float64 `json:"fee"`
	NeedMemo       bool    `json:"needMemo"`
	MinAmount      float64 `json:"minAmount"`
	MaxAmount      float64 `json:"maxAmount"`
}

// GetOffRampCurrencyPairsRequest filters are all optional and matched
// case-insensitively by the server.
type GetOffRampCurrencyPairsRequest struct {
	FiatCurrency   string `json:"fiatCurrency,omitempty"`
	CryptoCurrency string `json:"cryptoCurrency,omitempty"`
	Network        string `json:"network,omitempty"`
}

// GetOffRampCurrencyPairs lists the fiat/crypto/network combinations that
// can be sold via off-ramp, GET /v3/off-ramp/currency-pairs.
func (c *Client) GetOffRampCurrencyPairs(
	ctx context.Context,
	userEmail string,
	req *GetOffRampCurrencyPairsRequest,
) (*Response[[]CurrencyPair], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/off-ramp/currency-pairs",
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
