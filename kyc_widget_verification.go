package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetKycWidgetVerificationResponse is the caller's most recent
// widget-assisted KYC attempt, if any. The spec does not expand
// latestAttempt's inner fields, so it is left as a raw JSON value — confirm
// its shape against the live OpenAPI schema before relying on its contents.
type GetKycWidgetVerificationResponse struct {
	ProviderExternalUserID string          `json:"providerExternalUserId"`
	LatestAttempt          json.RawMessage `json:"latestAttempt,omitempty"`
}

// GetKycWidgetVerification returns the caller's most recent widget-assisted
// KYC attempt, GET /v3/kyc/widget-verifications.
func (c *Client) GetKycWidgetVerification(
	ctx context.Context,
	userEmail string,
) (*Response[GetKycWidgetVerificationResponse], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/kyc/widget-verifications",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[GetKycWidgetVerificationResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}

type CreateKycWidgetVerificationResponse struct {
	AttemptNumber          int    `json:"attemptNumber"`
	ProviderExternalUserID string `json:"providerExternalUserId"`
}

// Error code specific to creating a widget verification attempt.
const CodeWidgetKycMaxAttemptsReached = "WIDGET_KYC_MAX_ATTEMPTS_REACHED"

// CreateKycWidgetVerification starts a new widget-assisted KYC attempt for
// userEmail, POST /v3/kyc/widget-verifications.
func (c *Client) CreateKycWidgetVerification(
	ctx context.Context,
	userEmail string,
) (*Response[CreateKycWidgetVerificationResponse], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/kyc/widget-verifications",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[CreateKycWidgetVerificationResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
