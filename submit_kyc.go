package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type IDType string

const (
	IDTypeIDCard   IDType = "ID_CARD"
	IDTypePassport IDType = "PASSPORT"
)

type SubmitKycRequest struct {
	UserEmail        string `json:"userEmail"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	DateOfBirth      string `json:"dateOfBirth"`
	Gender           Gender `json:"gender"`
	Nationality      string `json:"nationality"`
	DocumentType     IDType `json:"type"`
	NationalID       string `json:"nationalId"`
	IssueDate        string `json:"issueDate"`
	ExpiryDate       string `json:"expiryDate"`
	AddressLine1     string `json:"addressLine1"`
	AddressLine2     string `json:"addressLine2"`
	City             string `json:"city"`
	State            string `json:"state"`
	ZipCode          string `json:"zipCode"`
	FrontIDImage     string `json:"frontIdImage"`
	BackIDImage      string `json:"backIdImage"`
	HoldIDImage      string `json:"holdIdImage"`
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	PhoneCountryCode string `json:"phoneCountryCode,omitempty"`
}

type SubmitKycResponse struct {
	NationalID string `json:"nationalId"`
	KycStatus  string `json:"kycStatus"`
	Signature  string `json:"signature"`
}

func (c *Client) SubmitKyc(ctx context.Context, req *SubmitKycRequest) (*Response[SubmitKycResponse], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s|%s",
		c.partnerCode,
		req.UserEmail,
		req.Nationality,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/user/submit-kyc",
		Params:      req,
		SigningData: []byte(signaturePayload),
		Header:      nil,
		Body:        nil,
		FullURL:     "",
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[SubmitKycResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
