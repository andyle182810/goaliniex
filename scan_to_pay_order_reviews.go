package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// ScanToPayOrderReviewRequest mirrors ScanToPayOrderRequest minus OrderID:
// it previews the same eligibility checks order creation runs, without
// moving funds.
type ScanToPayOrderReviewRequest struct {
	CryptoCurrency Currency     `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	BankCode       string       `json:"bankCode,omitempty"`
	AccountNumber  string       `json:"accountNumber,omitempty"`
	Description    string       `json:"description"`
	AccountName    string       `json:"accountName,omitempty"`
	QRToken        string       `json:"qrToken,omitempty"`
	QuoteID        string       `json:"quoteId,omitempty"`
}

// ScanToPayOrderReview is the quote-preview that would apply if the
// reviewed order were created now.
type ScanToPayOrderReview struct {
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Fee            float64      `json:"fee"`
	CryptoCurrency Currency     `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
}

// CreateScanToPayOrderReview validates a scan-to-pay order before creating
// it, without moving any funds, POST /v3/scan-to-pay/order-reviews.
func (c *Client) CreateScanToPayOrderReview(
	ctx context.Context,
	userEmail string,
	req *ScanToPayOrderReviewRequest,
) (*Response[ScanToPayOrderReview], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/scan-to-pay/order-reviews",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[ScanToPayOrderReview])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
