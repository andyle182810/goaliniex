package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type AccountResolutionRequest struct {
	FiatCurrency  FiatCurrency `json:"fiatCurrency"`
	BankCode      string       `json:"bankCode"`
	AccountNumber string       `json:"accountNumber"`
}

type AccountResolutionResponse struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

// Error code specific to account resolution.
const CodeFiatCurrencyNotSupported = "FIAT_CURRENCY_NOT_SUPPORTED"

// ResolveAccount resolves a bank account number to its registered
// account-holder name. Only available for VND, NGN, GHS, ZAR and KES,
// POST /v3/account-resolutions.
func (c *Client) ResolveAccount(
	ctx context.Context,
	userEmail string,
	req *AccountResolutionRequest,
) (*Response[AccountResolutionResponse], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/account-resolutions",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[AccountResolutionResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
