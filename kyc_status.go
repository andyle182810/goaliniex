package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// KycStatus is the user's KYC verification status. Note: this is spelled
// PROGRESSING (not PROCESSING, which is a distinct OrderStatus value).
type KycStatus string

const (
	KycStatusNone        KycStatus = "NONE"
	KycStatusProgressing KycStatus = "PROGRESSING"
	KycStatusVerified    KycStatus = "VERIFIED"
	KycStatusRejected    KycStatus = "REJECTED"
)

type GetKycStatusResponse struct {
	Status KycStatus `json:"status"`
}

// GetKycStatus returns userEmail's current KYC status, GET /v3/kyc/status.
func (c *Client) GetKycStatus(ctx context.Context, userEmail string) (*Response[GetKycStatusResponse], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/kyc/status",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[GetKycStatusResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
