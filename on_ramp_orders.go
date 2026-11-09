package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type OnRampOrderRequest struct {
	// OrderID is the caller-assigned idempotency key; reusing one is rejected.
	OrderID        string       `json:"orderId"`
	CryptoCurrency string       `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	Network        string       `json:"network"`
	WalletAddress  string       `json:"walletAddress"`
	CryptoAmount   float64      `json:"cryptoAmount,omitempty"`
	FiatAmount     float64      `json:"fiatAmount,omitempty"`
	Memo           string       `json:"memo,omitempty"`
	QuoteID        string       `json:"quoteId,omitempty"`
}

// OnRampPayment is the fiat payment instructions the caller must follow to
// settle an on-ramp order.
type OnRampPayment struct {
	BankName          string  `json:"bankName"`
	BankAccountName   string  `json:"bankAccountName"`
	BankAccountNumber string  `json:"bankAccountNumber"`
	Content           string  `json:"content"`
	Amount            float64 `json:"amount"`
	QRUrl             string  `json:"qrUrl"`
}

// OnRampDelivery is the on-chain crypto delivery for an on-ramp order.
type OnRampDelivery struct {
	CryptoCurrency string  `json:"cryptoCurrency"`
	Network        string  `json:"network"`
	Address        string  `json:"address"`
	Amount         float64 `json:"amount"`
	TxHash         string  `json:"txHash"`
}

// OnRampOrder is an on-ramp order, returned by both order creation and
// order-details lookups.
type OnRampOrder struct {
	OrderID      string         `json:"orderId"`
	Status       OrderStatus    `json:"status"`
	FiatCurrency FiatCurrency   `json:"fiatCurrency"`
	CreatedAt    string         `json:"createdAt"`
	ExpiresAt    string         `json:"expiresAt"`
	Payment      OnRampPayment  `json:"payment"`
	Delivery     OnRampDelivery `json:"delivery"`
	Fees         Fees           `json:"fees"`
	Quote        *Quote         `json:"quote,omitempty"`
}

// CreateOnRampOrder creates an on-ramp order buying crypto with fiat,
// delivered on-chain to the caller's wallet address, POST /v3/on-ramp/orders.
func (c *Client) CreateOnRampOrder(
	ctx context.Context,
	userEmail string,
	req *OnRampOrderRequest,
) (*Response[OnRampOrder], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/on-ramp/orders",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[OnRampOrder])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
