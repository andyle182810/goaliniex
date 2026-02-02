package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetWalletBalanceRequest struct {
	Currency Currency `json:"currency"`
}

type WalletBalance struct {
	Balance   float64  `json:"balance"`
	Currency  Currency `json:"currency"`
	Signature string   `json:"signature"`
}

func (c *Client) GetWalletBalance(ctx context.Context, req *GetWalletBalanceRequest) (*Response[WalletBalance], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s",
		c.partnerCode,
		req.Currency,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/wallet/balance",
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

	response := new(Response[WalletBalance])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
