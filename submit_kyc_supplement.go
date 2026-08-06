package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SubmitKycSupplementRequest struct {
	UserEmail        string `json:"userEmail"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	DateOfBirth      string `json:"dateOfBirth,omitempty"`
	Gender           Gender `json:"gender,omitempty"`
	Nationality      string `json:"nationality,omitempty"`
	DocumentType     IDType `json:"type,omitempty"`
	NationalID       string `json:"nationalId,omitempty"`
	IssueDate        string `json:"issueDate,omitempty"`
	ExpiryDate       string `json:"expiryDate,omitempty"`
	AddressLine1     string `json:"addressLine1,omitempty"`
	AddressLine2     string `json:"addressLine2,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	ZipCode          string `json:"zipCode,omitempty"`
	FrontIDImage     string `json:"frontIdImage,omitempty"`
	BackIDImage      string `json:"backIdImage,omitempty"`
	HoldIDImage      string `json:"holdIdImage,omitempty"`
	KycReport        string `json:"kycReport,omitempty"`
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	PhoneCountryCode string `json:"phoneCountryCode,omitempty"`
}

type SubmitKycSupplementResponse struct {
	ID                  int    `json:"id"`
	KycStatus           string `json:"kycStatus"`
	KycSupplementStatus string `json:"kycSupplementStatus"`
	Signature           string `json:"signature"`
}

func (c *Client) SubmitKycSupplement(
	ctx context.Context,
	req *SubmitKycSupplementRequest,
) (*Response[SubmitKycSupplementResponse], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s",
		c.partnerCode,
		req.UserEmail,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/user/submit-kyc-supplement",
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

	response := new(Response[SubmitKycSupplementResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
