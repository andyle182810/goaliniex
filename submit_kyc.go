package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// Gender is MALE/FEMALE per Alix v3 (not v2's lowercase male/female).
type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type IDType string

const (
	IDTypeIDCard          IDType = "ID_CARD"
	IDTypePassport        IDType = "PASSPORT"
	IDTypeDriversLicense  IDType = "DRIVERS_LICENSE"
	IDTypeResidencePermit IDType = "RESIDENCE_PERMIT"
)

// SubmitKycAddress is the optional residential address group for a KYC
// submission/supplement. Line1, City, State and ZipCode are required
// together; Line2 is optional.
type SubmitKycAddress struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2,omitempty"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
}

// SubmitKycPhone is the optional contact phone group for a KYC
// submission/supplement.
type SubmitKycPhone struct {
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

type SubmitKycRequest struct {
	FirstName    string            `json:"firstName"`
	LastName     string            `json:"lastName"`
	Nationality  string            `json:"nationality"`
	DocumentType IDType            `json:"documentType"`
	DateOfBirth  string            `json:"dateOfBirth"`
	Gender       Gender            `json:"gender,omitempty"`
	NationalID   string            `json:"nationalId"`
	IssueDate    string            `json:"issueDate,omitempty"`
	ExpiryDate   string            `json:"expiryDate"`
	Address      *SubmitKycAddress `json:"address,omitempty"`
	Phone        *SubmitKycPhone   `json:"phone,omitempty"`
	FrontImage   string            `json:"frontImage"`
	// BackImage is required unless DocumentType is IDTypePassport.
	BackImage      string `json:"backImage,omitempty"`
	SelfieImage    string `json:"selfieImage"`
	ReportDocument string `json:"reportDocument"`
}

type SubmitKycResponse struct {
	Status KycStatus `json:"status"`
}

// SubmitKyc submits userEmail's identity information for KYC verification,
// POST /v3/kyc/submissions.
func (c *Client) SubmitKyc(
	ctx context.Context,
	userEmail string,
	req *SubmitKycRequest,
) (*Response[SubmitKycResponse], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/kyc/submissions",
		Params:    req,
		UserEmail: userEmail,
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
