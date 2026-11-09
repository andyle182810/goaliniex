package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// DepositAddress is an on-chain address for topping up the partner's crypto
// fund on a given network.
type DepositAddress struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

// PartnerFund is the partner's crypto fund balance and its deposit
// addresses.
type PartnerFund struct {
	CryptoCurrency   string           `json:"cryptoCurrency"`
	Balance          float64          `json:"balance"`
	DepositAddresses []DepositAddress `json:"depositAddresses"`
}

// GetPartnerFund returns the partner's crypto fund balance and deposit
// addresses, GET /v3/partner/fund. cryptoCurrency defaults to USDT when
// empty.
func (c *Client) GetPartnerFund(ctx context.Context, cryptoCurrency string) (*Response[PartnerFund], error) {
	params := map[string]any{}
	if cryptoCurrency != "" {
		params["cryptoCurrency"] = cryptoCurrency
	}

	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/v3/partner/fund",
		Params:   params,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[PartnerFund])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
