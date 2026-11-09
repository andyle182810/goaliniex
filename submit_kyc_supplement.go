package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// SubmitKycSupplementRequest fields are all optional: fields sent replace
// the stored values, the rest stay as they are.
type SubmitKycSupplementRequest struct {
	FirstName      string            `json:"firstName,omitempty"`
	LastName       string            `json:"lastName,omitempty"`
	Nationality    string            `json:"nationality,omitempty"`
	DocumentType   IDType            `json:"documentType,omitempty"`
	DateOfBirth    string            `json:"dateOfBirth,omitempty"`
	Gender         Gender            `json:"gender,omitempty"`
	NationalID     string            `json:"nationalId,omitempty"`
	IssueDate      string            `json:"issueDate,omitempty"`
	ExpiryDate     string            `json:"expiryDate,omitempty"`
	Address        *SubmitKycAddress `json:"address,omitempty"`
	Phone          *SubmitKycPhone   `json:"phone,omitempty"`
	FrontImage     string            `json:"frontImage,omitempty"`
	BackImage      string            `json:"backImage,omitempty"`
	SelfieImage    string            `json:"selfieImage,omitempty"`
	ReportDocument string            `json:"reportDocument,omitempty"`
}

// SubmitKycSupplementResponse carries only Status: a supplement never
// changes the user's KYC status.
type SubmitKycSupplementResponse struct {
	Status KycStatus `json:"status"`
}

// SubmitKycSupplement adds to or corrects userEmail's submitted KYC profile,
// POST /v3/kyc/supplements.
func (c *Client) SubmitKycSupplement(
	ctx context.Context,
	userEmail string,
	req *SubmitKycSupplementRequest,
) (*Response[SubmitKycSupplementResponse], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/kyc/supplements",
		Params:    req,
		UserEmail: userEmail,
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
