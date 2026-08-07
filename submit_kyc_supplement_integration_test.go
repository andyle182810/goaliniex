package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_SubmitKycSupplement_BackfillKycReport(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        email,
		FirstName:        "",
		LastName:         "",
		DateOfBirth:      "",
		Gender:           "",
		Nationality:      "",
		DocumentType:     "",
		NationalID:       "",
		IssueDate:        "",
		ExpiryDate:       "",
		AddressLine1:     "",
		AddressLine2:     "",
		City:             "",
		State:            "",
		ZipCode:          "",
		FrontIDImage:     "",
		BackIDImage:      "",
		HoldIDImage:      "",
		KycReport:        getTestPdfDataURI(),
		PhoneNumber:      "",
		PhoneCountryCode: "",
	}

	resp, err := client.SubmitKycSupplement(ctx, req)
	if err != nil {
		t.Fatalf("SubmitKycSupplement returned error: %v", err)
	}

	t.Logf("SubmitKycSupplement response: success=%v, message=%s", resp.Success, resp.Message)

	if resp.Success && resp.Data != nil {
		t.Logf("  ID: %d", resp.Data.ID)
		t.Logf("  KycStatus: %s", resp.Data.KycStatus)
		t.Logf("  KycSupplementStatus: %s", resp.Data.KycSupplementStatus)
	}
}

func TestIntegration_SubmitKycSupplement_UnknownUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        generateRandomGmail(t),
		FirstName:        "",
		LastName:         "",
		DateOfBirth:      "",
		Gender:           "",
		Nationality:      "",
		DocumentType:     "",
		NationalID:       "",
		IssueDate:        "",
		ExpiryDate:       "",
		AddressLine1:     "",
		AddressLine2:     "",
		City:             "",
		State:            "",
		ZipCode:          "",
		FrontIDImage:     "",
		BackIDImage:      "",
		HoldIDImage:      "",
		KycReport:        getTestPdfDataURI(),
		PhoneNumber:      "",
		PhoneCountryCode: "",
	}

	resp, err := client.SubmitKycSupplement(ctx, req)
	if err != nil {
		t.Logf("SubmitKycSupplement returned error for unknown user: %v", err)

		return
	}

	t.Logf("Response for unknown user: success=%v, message=%s", resp.Success, resp.Message)
}
