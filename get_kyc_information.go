package goaliniex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type KycStatus string

const (
	KycStatusNone       KycStatus = "NONE"
	KycStatusProcessing KycStatus = "PROCESSING"
	KycStatusVerified   KycStatus = "VERIFIED"
	KycStatusRejected   KycStatus = "REJECTED"
)

type KycInformationRequest struct {
	UserEmail string `json:"userEmail"`
}

type KycInformation struct {
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	DateOfBirth      string    `json:"dateOfBirth"`
	Gender           string    `json:"gender"`
	Nationality      string    `json:"nationality"`
	IDType           string    `json:"idType"`
	NationalID       string    `json:"nationalId"`
	IssueDate        string    `json:"issueDate"`
	ExpiryDate       string    `json:"expiryDate"`
	Address          string    `json:"address"`
	FrontIDImage     string    `json:"frontIdImage"`
	BackIDImage      string    `json:"backIdImage"`
	HoldIDImage      string    `json:"holdIdImage"`
	PhoneNumber      string    `json:"phoneNumber"`
	PhoneCountryCode string    `json:"phoneCountryCode"`
	KycStatus        KycStatus `json:"kycStatus"`
	RejectReason     string    `json:"rejectReason"`
}

func (c *Client) GetKycInformation(ctx context.Context, req *KycInformationRequest) (*Response[KycInformation], error) {
	signaturePayload := fmt.Sprintf(
		"%s|%s|%s",
		c.partnerCode,
		req.UserEmail,
		c.secretKey,
	)

	apiRequest := request{
		Method:      http.MethodPost,
		Endpoint:    "/api/v2/user/get-kyc-information",
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

	response := new(Response[KycInformation])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
