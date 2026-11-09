package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// CashoutRecipient is the bank account a cashout order pays out to. Unlike
// scan-to-pay's top-level BankCode/AccountNumber, AccountName is required
// here.
type CashoutRecipient struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Description   string `json:"description"`
	QRType        string `json:"qrType,omitempty"`
}

type CashoutOrderRequest struct {
	// OrderID is the caller-assigned idempotency key; reusing one is rejected.
	OrderID        string           `json:"orderId"`
	CryptoCurrency Currency         `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency     `json:"fiatCurrency"`
	FiatAmount     float64          `json:"fiatAmount"`
	Recipient      CashoutRecipient `json:"recipient"`
	QuoteID        string           `json:"quoteId,omitempty"`
}

// CashoutOrder is a cashout order, returned by both order creation and
// order-details lookups.
type CashoutOrder struct {
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

// CreateCashoutOrder creates a cashout order paying fiat to a recipient bank
// account, debiting the partner's crypto fund, POST /v3/cashout/orders.
func (c *Client) CreateCashoutOrder(
	ctx context.Context,
	userEmail string,
	req *CashoutOrderRequest,
) (*Response[CashoutOrder], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/cashout/orders",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[CashoutOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
