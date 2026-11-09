package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// KycAddress is the user's residential address. Fields are empty strings
// when not on file.
type KycAddress struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
}

// KycPhone is the user's contact phone number. Fields are empty strings
// when not on file.
type KycPhone struct {
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

// KycProfile is the user's KYC profile, available in any state once a
// submission exists. Unlike v2's KycInformation, it does not return the
// document images (front/back/selfie/report) and Address/Phone are
// structured rather than flat strings.
type KycProfile struct {
	Status       KycStatus  `json:"status"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	DateOfBirth  string     `json:"dateOfBirth"`
	Gender       Gender     `json:"gender"`
	Nationality  string     `json:"nationality"`
	DocumentType IDType     `json:"documentType"`
	NationalID   string     `json:"nationalId"`
	IssueDate    string     `json:"issueDate"`
	ExpiryDate   string     `json:"expiryDate"`
	Address      KycAddress `json:"address"`
	Phone        KycPhone   `json:"phone"`
	RejectReason string     `json:"rejectReason"`
}

// GetKycProfile returns userEmail's KYC profile, GET /v3/kyc/profile.
func (c *Client) GetKycProfile(ctx context.Context, userEmail string) (*Response[KycProfile], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/kyc/profile",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[KycProfile])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
