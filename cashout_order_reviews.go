package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// CashoutOrderReviewRequest mirrors CashoutOrderRequest minus OrderID: it
// previews the same eligibility checks order creation runs, without moving
// funds.
type CashoutOrderReviewRequest struct {
	CryptoCurrency Currency         `json:"cryptoCurrency,omitempty"`
	FiatCurrency   FiatCurrency     `json:"fiatCurrency"`
	FiatAmount     float64          `json:"fiatAmount"`
	Recipient      CashoutRecipient `json:"recipient"`
	QuoteID        string           `json:"quoteId,omitempty"`
}

// CashoutOrderReview is the quote-preview that would apply if the reviewed
// order were created now.
type CashoutOrderReview struct {
	FiatCurrency   FiatCurrency `json:"fiatCurrency"`
	FiatAmount     float64      `json:"fiatAmount"`
	Fee            float64      `json:"fee"`
	CryptoCurrency Currency     `json:"cryptoCurrency"`
	CryptoAmount   float64      `json:"cryptoAmount"`
}

// CreateCashoutOrderReview validates a cashout order before creating it,
// without moving any funds, POST /v3/cashout/order-reviews.
func (c *Client) CreateCashoutOrderReview(
	ctx context.Context,
	userEmail string,
	req *CashoutOrderReviewRequest,
) (*Response[CashoutOrderReview], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/cashout/order-reviews",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[CashoutOrderReview])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
