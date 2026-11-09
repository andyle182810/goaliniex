package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex/v2"
)

func TestIntegration_SubmitKyc_Valid(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := generateRandomGmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		FirstName:      "Test",
		LastName:       "User",
		Nationality:    "VN",
		DocumentType:   goaliniex.IDTypePassport,
		DateOfBirth:    "1990-01-01",
		Gender:         goaliniex.GenderMale,
		NationalID:     generateTestSSN(t),
		ExpiryDate:     "2035-01-01",
		FrontImage:     getTestImageDataURI(),
		SelfieImage:    getTestImageDataURI(),
		ReportDocument: getTestPdfDataURI(),
	}

	resp, err := client.SubmitKyc(ctx, email, req)
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	t.Logf("SubmitKyc response: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("KYC status after submission: %s", resp.Data.Status)
	}
}

func TestIntegration_GetKycStatus(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.GetKycStatus(ctx, email)
	if err != nil {
		t.Fatalf("GetKycStatus returned error: %v", err)
	}

	t.Logf("GetKycStatus: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("status=%s", resp.Data.Status)
	}
}

func TestIntegration_GetKycProfile(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.GetKycProfile(ctx, email)
	if err != nil {
		t.Fatalf("GetKycProfile returned error: %v", err)
	}

	t.Logf("GetKycProfile: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("firstName=%s lastName=%s status=%s", resp.Data.FirstName, resp.Data.LastName, resp.Data.Status)
	}
}

func TestIntegration_SubmitKycSupplement(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.SubmitKycSupplement(ctx, email, &goaliniex.SubmitKycSupplementRequest{
		Phone: &goaliniex.SubmitKycPhone{Number: "123456789", CountryCode: "84"},
	})
	if err != nil {
		t.Fatalf("SubmitKycSupplement returned error: %v", err)
	}

	t.Logf("SubmitKycSupplement: code=%s message=%s", resp.Code, resp.Message)
}
