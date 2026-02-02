package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Currency string

const (
	CurrencyUSDT Currency = "USDT"
	CurrencyETH  Currency = "ETH"
	CurrencyBTC  Currency = "BTC"
)

type FiatCurrency string

const (
	FiatCurrencyVND FiatCurrency = "VND"
	FiatCurrencyPHP FiatCurrency = "PHP"
	FiatCurrencyTHB FiatCurrency = "THB"
	FiatCurrencyGEL FiatCurrency = "GEL"
	FiatCurrencyBRL FiatCurrency = "BRL"
	FiatCurrencyARS FiatCurrency = "ARS"
	FiatCurrencyPEN FiatCurrency = "PEN"
	FiatCurrencyNGN FiatCurrency = "NGN"
)

type OrderStatus string

const (
	OrderStatusAwaitingPayment         OrderStatus = "AWAITING_PAYMENT"
	OrderStatusPaymentCompleted        OrderStatus = "PAYMENT_COMPLETED"
	OrderStatusProcessingTokenTransfer OrderStatus = "PROCESSING_TOKEN_TRANSFER"
	OrderStatusSuccess                 OrderStatus = "SUCCESS"
	OrderStatusError                   OrderStatus = "ERROR"
	OrderStatusFail                    OrderStatus = "FAIL"
)

type CreateOrderRequest struct {
	Currency          Currency     `json:"currency"`
	FiatAmount        float64      `json:"fiatAmount"`
	FiatCurrency      FiatCurrency `json:"fiatCurrency"`
	BankCode          string       `json:"bankCode"`
	BankAccountNumber string       `json:"bankAccountNumber"`
	ExternalOrderID   string       `json:"externalOrderId"`
	WebhookSecretKey  string       `json:"webhookSecretKey"`
	UserEmail         string       `json:"userEmail"`
	UserKYCVerified   bool         `json:"userKycVerified"`
	Content           string       `json:"content"`
	ExtendInfo        any          `json:"extendInfo"`
}

type CreateOrderResponse struct {
	ExternalOrderID string        `json:"externalOrderId"`
	Type            string        `json:"type"`
	FiatAmount      float64       `json:"fiatAmount"`
	PaidAmount      float64       `json:"paidAmount"`
	TokenTransfer   TokenTransfer `json:"tokenTransfer"`
	BankTransfer    BankTransfer  `json:"bankTransfer"`
	Fees            Fees          `json:"fees"`
	Status          OrderStatus   `json:"status"`
	Descriptions    string        `json:"descriptions"`
	CreatedAt       string        `json:"createdAt"`
	ExpiresAt       string        `json:"expiresAt"`
	Signature       string        `json:"signature"`
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Response[CreateOrderResponse], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%s|%s|%s",
		c.partnerCode,
		req.ExternalOrderID,
		req.Currency,
		strconv.FormatFloat(req.FiatAmount, 'f', -1, 64),
		req.BankCode,
		req.BankAccountNumber,
		req.Content,
		req.UserEmail,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/orders/create-sell-order",
		Params:      req,
		SigningData: []byte(signaturePayload),
		Header:      nil,
		Body:        nil,
		FullURL:     "",
		Public:      false,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[CreateOrderResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
