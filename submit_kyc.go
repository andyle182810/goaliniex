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
	IDTypeIDCard          IDType = "ID_CARD"
	IDTypePassport        IDType = "PASSPORT"
	IDTypeDriversLicense  IDType = "DRIVERS_LICENSE"
	IDTypeResidencePermit IDType = "RESIDENCE_PERMIT"
)

type SubmitKycRequest struct {
	UserEmail    string `json:"userEmail"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	DateOfBirth  string `json:"dateOfBirth"`
	Gender       Gender `json:"gender,omitempty"`
	Nationality  string `json:"nationality"`
	DocumentType IDType `json:"type"`
	NationalID   string `json:"nationalId"`
	IssueDate    string `json:"issueDate,omitempty"`
	ExpiryDate   string `json:"expiryDate"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	ZipCode      string `json:"zipCode,omitempty"`
	FrontIDImage string `json:"frontIdImage"`
	// BackIDImage is required unless DocumentType is IDTypePassport.
	BackIDImage      string `json:"backIdImage,omitempty"`
	HoldIDImage      string `json:"holdIdImage"`
	KycReport        string `json:"kycReport"`
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	PhoneCountryCode string `json:"phoneCountryCode,omitempty"`
}

type SubmitKycResponse struct {
	ID        int    `json:"id"`
	KycStatus string `json:"kycStatus"`
	Signature string `json:"signature"`
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
		Public:      false,
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
