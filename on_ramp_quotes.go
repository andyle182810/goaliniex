package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// OnRampQuoteRequest previews buying crypto with fiat. CryptoAmount and
// FiatAmount are mutually exclusive: set exactly one.
type OnRampQuoteRequest struct {
	CryptoCurrency string       `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	Network        string       `json:"network"`
	CryptoAmount   float64      `json:"cryptoAmount,omitempty"`
	FiatAmount     float64      `json:"fiatAmount,omitempty"`
	Lock           bool         `json:"lock,omitempty"`
}

// OnRampQuote previews (or, with Lock, locks) the fiat cost and fee for
// buying crypto.
type OnRampQuote struct {
	QuoteID        string       `json:"quoteId,omitempty"`
	Price          float64      `json:"price,omitempty"`
	ExpiresAt      string       `json:"expiresAt,omitempty"`
	CryptoCurrency string       `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Fees           Fees         `json:"fees"`
}

// CreateOnRampQuote previews the fiat cost and fee for buying crypto,
// POST /v3/on-ramp/quotes.
func (c *Client) CreateOnRampQuote(
	ctx context.Context,
	userEmail string,
	req *OnRampQuoteRequest,
) (*Response[OnRampQuote], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/on-ramp/quotes",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[OnRampQuote])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
