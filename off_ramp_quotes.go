package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// OffRampQuoteRequest previews selling crypto for fiat. CryptoAmount and
// FiatAmount are mutually exclusive: set exactly one.
type OffRampQuoteRequest struct {
	CryptoCurrency string       `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	Network        string       `json:"network"`
	CryptoAmount   float64      `json:"cryptoAmount,omitempty"`
	FiatAmount     float64      `json:"fiatAmount,omitempty"`
	Lock           bool         `json:"lock,omitempty"`
}

// OffRampQuote previews (or, with Lock, locks) the net fiat payout and fee
// for selling crypto.
type OffRampQuote struct {
	QuoteID        string       `json:"quoteId,omitempty"`
	Price          float64      `json:"price,omitempty"`
	ExpiresAt      string       `json:"expiresAt,omitempty"`
	CryptoCurrency string       `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Fees           Fees         `json:"fees"`
}

// CreateOffRampQuote previews the net fiat payout and fee for selling
// crypto, POST /v3/off-ramp/quotes.
func (c *Client) CreateOffRampQuote(
	ctx context.Context,
	userEmail string,
	req *OffRampQuoteRequest,
) (*Response[OffRampQuote], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/off-ramp/quotes",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[OffRampQuote])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
