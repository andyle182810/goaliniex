package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type ScanToPayQuoteRequest struct {
	CryptoCurrency Currency     `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Lock           bool         `json:"lock,omitempty"`
}

// ScanToPayQuote previews (or, with Lock, locks) the crypto cost and fee for
// a scan-to-pay payout.
type ScanToPayQuote struct {
	QuoteID        string       `json:"quoteId,omitempty"`
	Price          float64      `json:"price,omitempty"`
	ExpiresAt      string       `json:"expiresAt,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Fee            float64      `json:"fee"`
	CryptoCurrency Currency     `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
	MinAmount      *float64     `json:"minAmount"`
	MaxAmount      *float64     `json:"maxAmount"`
}

// CreateScanToPayQuote previews the crypto cost and fee for a scan-to-pay
// payout, POST /v3/scan-to-pay/quotes.
func (c *Client) CreateScanToPayQuote(
	ctx context.Context,
	userEmail string,
	req *ScanToPayQuoteRequest,
) (*Response[ScanToPayQuote], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/scan-to-pay/quotes",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[ScanToPayQuote])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
