package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// PartnerService is one product the partner can use, and the fiat
// currencies enabled for it.
type PartnerService struct {
	Service        string   `json:"service"`
	FiatCurrencies []string `json:"fiatCurrencies"`
}

// GetPartnerServices lists the services the partner can use and the fiat
// currencies enabled for each, GET /v3/partner/services.
func (c *Client) GetPartnerServices(ctx context.Context) (*Response[[]PartnerService], error) {
	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/v3/partner/services",
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]PartnerService])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
