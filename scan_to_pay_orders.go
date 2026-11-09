package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type ScanToPayOrderRequest struct {
	// OrderID is the caller-assigned idempotency key; reusing one is rejected.
	OrderID        string       `json:"orderId"`
	CryptoCurrency Currency     `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	// BankCode and AccountNumber are omitted when the QR's payeeAddressing
	// is PROXY or SCHEME_TOKEN.
	BankCode      string `json:"bankCode,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
	Description   string `json:"description"`
	AccountName   string `json:"accountName,omitempty"`
	// QRToken is the opaque token from DecodeQR, passed back verbatim.
	QRToken string `json:"qrToken,omitempty"`
	// QuoteID locks the order to a previously-locked quote's rate and fees.
	QuoteID string `json:"quoteId,omitempty"`
}

// ScanToPayOrder is a scan-to-pay order, returned by both order creation and
// order-details lookups.
type ScanToPayOrder struct {
	OrderID        string       `json:"orderId"`
	Status         OrderStatus  `json:"status"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	CryptoCurrency Currency     `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
	CreatedAt      string       `json:"createdAt"`
	Payout         Payout       `json:"payout"`
	Fee            float64      `json:"fee"`
	Quote          *Quote       `json:"quote,omitempty"`
}

// CreateScanToPayOrder creates a scan-to-pay order paying fiat to a bank
// account or QR-identified recipient, debiting the partner's crypto fund,
// POST /v3/scan-to-pay/orders.
func (c *Client) CreateScanToPayOrder(
	ctx context.Context,
	userEmail string,
	req *ScanToPayOrderRequest,
) (*Response[ScanToPayOrder], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/scan-to-pay/orders",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[ScanToPayOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
