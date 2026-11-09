package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// OffRampRecipient is the fiat payout destination for an off-ramp order.
// AccountName is optional only for VND; required for every other currency.
type OffRampRecipient struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName,omitempty"`
}

type OffRampOrderRequest struct {
	// OrderID is the caller-assigned idempotency key; reusing one is rejected.
	OrderID        string           `json:"orderId"`
	CryptoCurrency string           `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency     `json:"fiatCurrency"`
	Network        string           `json:"network"`
	CryptoAmount   float64          `json:"cryptoAmount,omitempty"`
	FiatAmount     float64          `json:"fiatAmount,omitempty"`
	Recipient      OffRampRecipient `json:"recipient"`
	QuoteID        string           `json:"quoteId,omitempty"`
}

// OffRampDeposit is the on-chain crypto deposit the caller must send to
// settle an off-ramp order.
type OffRampDeposit struct {
	CryptoCurrency string  `json:"cryptoCurrency"`
	Network        string  `json:"network"`
	Address        string  `json:"address"`
	Amount         float64 `json:"amount"`
	Memo           string  `json:"memo"`
}

// OffRampOrder is an off-ramp order, returned by both order creation and
// order-details lookups.
type OffRampOrder struct {
	OrderID      string         `json:"orderId"`
	Status       OrderStatus    `json:"status"`
	FiatCurrency FiatCurrency   `json:"fiatCurrency"`
	CreatedAt    string         `json:"createdAt"`
	ExpiresAt    string         `json:"expiresAt"`
	Deposit      OffRampDeposit `json:"deposit"`
	Payout       Payout         `json:"payout"`
	Fees         Fees           `json:"fees"`
	Quote        *Quote         `json:"quote,omitempty"`
}

// CreateOffRampOrder creates an off-ramp order selling crypto for fiat, paid
// out to a recipient bank account, POST /v3/off-ramp/orders.
func (c *Client) CreateOffRampOrder(
	ctx context.Context,
	userEmail string,
	req *OffRampOrderRequest,
) (*Response[OffRampOrder], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/off-ramp/orders",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[OffRampOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
