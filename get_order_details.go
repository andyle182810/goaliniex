package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetOrderDetailsRequest struct {
	ExternalOrderID string `json:"externalOrderId"`
}

type TokenTransfer struct {
	Currency      Currency `json:"currency"`
	Network       string   `json:"network"`
	Price         float64  `json:"price"`
	Amount        float64  `json:"amount"`
	WalletAddress string   `json:"walletAddress"`
	TxHash        string   `json:"txHash"`
}

type BankTransfer struct {
	BankCode          string `json:"bankCode"`
	BankName          string `json:"bankName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankAccountName   string `json:"bankAccountName"`
	Content           string `json:"content"`
	QRCodeURL         string `json:"qrCodeUrl"`
}

type Fees struct {
	SystemFee     float64 `json:"systemFee"`
	ProcessingFee float64 `json:"processingFee"`
}

type OrderDetails struct {
	ExternalOrderID string        `json:"externalOrderId"`
	Type            string        `json:"type"`
	FiatAmount      float64       `json:"fiatAmount"`
	FiatCurrency    FiatCurrency  `json:"fiatCurrency"`
	PaidAmount      float64       `json:"paidAmount"`
	TokenTransfer   TokenTransfer `json:"tokenTransfer"`
	BankTransfer    BankTransfer  `json:"bankTransfer"`
	Fees            Fees          `json:"fees"`
	Status          OrderStatus   `json:"status"`
	Description     string        `json:"description"`
	CreatedAt       string        `json:"createdAt"`
	ExpiresAt       string        `json:"expiresAt"`
	Signature       string        `json:"signature"`
}

func (c *Client) GetOrderDetails(ctx context.Context, req *GetOrderDetailsRequest) (*Response[OrderDetails], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s",
		c.partnerCode,
		req.ExternalOrderID,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/orders/details",
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

	response := new(Response[OrderDetails])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
