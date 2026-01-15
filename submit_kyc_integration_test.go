package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_SubmitKyc_ValidRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "tvugiang@gmail.com",
		FirstName:        "John",
		LastName:         "Doe",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "123456789",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Main St",
		AddressLine2:     "Apt 4",
		City:             "New York",
		State:            "NY",
		ZipCode:          "10001",
		FrontIDImage:     "base64frontimage",
		BackIDImage:      "base64backimage",
		HoldIDImage:      "base64holdimage",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	t.Logf("SubmitKyc response: success=%v, message=%s", resp.Success, resp.Message)

	if resp.Success && resp.Data != nil {
		t.Logf("  NationalID: %s", resp.Data.NationalID)
		t.Logf("  KycStatus: %s", resp.Data.KycStatus)
		t.Logf("  Signature length: %d", len(resp.Data.Signature))
	}
}

func TestIntegration_SubmitKyc_WithIDCard(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "zavytran1409+3@gmail.com",
		FirstName:        "Jane",
		LastName:         "Smith",
		DateOfBirth:      "1985-05-15",
		Gender:           goaliniex.GenderFemale,
		Nationality:      "VN",
		DocumentType:     goaliniex.IDTypeIDCard,
		NationalID:       "987654321",
		IssueDate:        "2019-06-01",
		ExpiryDate:       "2029-06-01",
		AddressLine1:     "456 Oak Ave",
		AddressLine2:     "456 Oak Ave",
		City:             "Ho Chi Minh",
		State:            "HCM",
		ZipCode:          "70000",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "0987654321",
		PhoneCountryCode: "84",
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	t.Logf("SubmitKyc response: success=%v, message=%s", resp.Success, resp.Message)
}

func TestIntegration_SubmitKyc_EmptyEmail(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "",
		FirstName:        "Test",
		LastName:         "User",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "111222333",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Test St",
		AddressLine2:     "123 Test St",
		City:             "Test City",
		State:            "TS",
		ZipCode:          "12345",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Logf("SubmitKyc returned error for empty email: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted empty email, this might be a validation issue")
	} else {
		t.Logf("API correctly rejected empty email: %s", resp.Message)
	}
}

func TestIntegration_SubmitKyc_InvalidEmailFormat(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "invalid-email-format",
		FirstName:        "Test",
		LastName:         "User",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "111222333",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Test St",
		AddressLine2:     "123 Test St",
		City:             "Test City",
		State:            "TS",
		ZipCode:          "12345",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Logf("SubmitKyc returned error for invalid email: %v", err)

		return
	}

	t.Logf("Response success: %v, message: %s", resp.Success, resp.Message)
}

func TestIntegration_SubmitKyc_MissingRequiredFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{ //nolint:exhaustruct
		UserEmail: "test@example.com",
		FirstName: "Test",
		// Missing other required fields
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Logf("SubmitKyc returned error for missing fields: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted request with missing required fields")
	} else {
		t.Logf("API correctly rejected incomplete request: %s", resp.Message)
	}
}

func TestIntegration_SubmitKyc_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "test@example.com",
		FirstName:        "Test",
		LastName:         "User",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "123456789",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Test St",
		AddressLine2:     "123 Test St",
		City:             "Test City",
		State:            "TS",
		ZipCode:          "12345",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	_, err := client.SubmitKyc(ctx, req)
	if err == nil {
		t.Log("Request completed before timeout (fast network)")
	} else {
		if strings.Contains(err.Error(), "context") ||
			strings.Contains(err.Error(), "deadline") ||
			strings.Contains(err.Error(), "timeout") {
			t.Logf("Context timeout error as expected: %v", err)
		} else {
			t.Logf("Got error (may or may not be timeout related): %v", err)
		}
	}
}

func TestIntegration_SubmitKyc_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "test@example.com",
		FirstName:        "Test",
		LastName:         "User",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "123456789",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Test St",
		AddressLine2:     "123 Test St",
		City:             "Test City",
		State:            "TS",
		ZipCode:          "12345",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	cancel()

	_, err := client.SubmitKyc(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}

func TestIntegration_SubmitKyc_DifferentNationalities(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	nationalities := []string{
		"US",
		"VN",
		"JP",
		"UK",
	}

	for _, nationality := range nationalities {
		t.Run(nationality, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.SubmitKycRequest{
				UserEmail:        "test_" + nationality + "@example.com",
				FirstName:        "Test",
				LastName:         "User",
				DateOfBirth:      "1990-01-01",
				Gender:           goaliniex.GenderMale,
				Nationality:      nationality,
				DocumentType:     goaliniex.IDTypePassport,
				NationalID:       "123456789",
				IssueDate:        "2020-01-01",
				ExpiryDate:       "2030-01-01",
				AddressLine1:     "123 Test St",
				AddressLine2:     "123 Test St",
				City:             "Test City",
				State:            "TS",
				ZipCode:          "12345",
				FrontIDImage:     "base64front",
				BackIDImage:      "base64back",
				HoldIDImage:      "base64hold",
				PhoneNumber:      "1234567890",
				PhoneCountryCode: "1",
			}

			resp, err := client.SubmitKyc(ctx, req)
			if err != nil {
				t.Logf("SubmitKyc error for nationality %s: %v", nationality, err)

				return
			}

			t.Logf("Response for nationality %s: success=%v, message=%s",
				nationality, resp.Success, resp.Message)
		})
	}
}

func TestIntegration_SubmitKyc_DocumentTypes(t *testing.T) { //nolint:dupl
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	docTypes := []struct {
		name    string
		docType goaliniex.IDType
	}{
		{"Passport", goaliniex.IDTypePassport},
		{"IDCard", goaliniex.IDTypeIDCard},
	}

	for _, docType := range docTypes {
		t.Run(docType.name, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.SubmitKycRequest{
				UserEmail:        "test_" + docType.name + "@example.com",
				FirstName:        "Test",
				LastName:         "User",
				DateOfBirth:      "1990-01-01",
				Gender:           goaliniex.GenderMale,
				Nationality:      "US",
				DocumentType:     docType.docType,
				NationalID:       "123456789",
				IssueDate:        "2020-01-01",
				ExpiryDate:       "2030-01-01",
				AddressLine1:     "123 Test St",
				AddressLine2:     "123 Test St",
				City:             "Test City",
				State:            "TS",
				ZipCode:          "12345",
				FrontIDImage:     "base64front",
				BackIDImage:      "base64back",
				HoldIDImage:      "base64hold",
				PhoneNumber:      "1234567890",
				PhoneCountryCode: "1",
			}

			resp, err := client.SubmitKyc(ctx, req)
			if err != nil {
				t.Logf("SubmitKyc error for document type %s: %v", docType.name, err)

				return
			}

			t.Logf("Response for document type %s: success=%v, message=%s",
				docType.name, resp.Success, resp.Message)
		})
	}
}

func TestIntegration_SubmitKyc_GenderValues(t *testing.T) { //nolint:dupl
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	genders := []struct {
		name   string
		gender goaliniex.Gender
	}{
		{"Male", goaliniex.GenderMale},
		{"Female", goaliniex.GenderFemale},
	}

	for _, gender := range genders {
		t.Run(gender.name, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.SubmitKycRequest{
				UserEmail:        "test_" + gender.name + "@example.com",
				FirstName:        "Test",
				LastName:         "User",
				DateOfBirth:      "1990-01-01",
				Gender:           gender.gender,
				Nationality:      "US",
				DocumentType:     goaliniex.IDTypePassport,
				NationalID:       "123456789",
				IssueDate:        "2020-01-01",
				ExpiryDate:       "2030-01-01",
				AddressLine1:     "123 Test St",
				AddressLine2:     "123 Test St",
				City:             "Test City",
				State:            "TS",
				ZipCode:          "12345",
				FrontIDImage:     "base64front",
				BackIDImage:      "base64back",
				HoldIDImage:      "base64hold",
				PhoneNumber:      "1234567890",
				PhoneCountryCode: "1",
			}

			resp, err := client.SubmitKyc(ctx, req)
			if err != nil {
				t.Logf("SubmitKyc error for gender %s: %v", gender.name, err)

				return
			}

			t.Logf("Response for gender %s: success=%v, message=%s",
				gender.name, resp.Success, resp.Message)
		})
	}
}

func TestIntegration_SubmitKyc_LongTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "tvugiang@gmail.com",
		FirstName:        "John",
		LastName:         "Doe",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "123456789",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Main St",
		AddressLine2:     "123 Main St",
		City:             "New York",
		State:            "NY",
		ZipCode:          "10001",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	start := time.Now()
	resp, err := client.SubmitKyc(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	t.Logf("Request completed in %v", duration)
	t.Logf("Response: success=%v, message=%s", resp.Success, resp.Message)
}

func TestIntegration_SubmitKyc_SpecialCharactersInAddress(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "test@example.com",
		FirstName:        "Test",
		LastName:         "User",
		DateOfBirth:      "1990-01-01",
		Gender:           goaliniex.GenderMale,
		Nationality:      "US",
		DocumentType:     goaliniex.IDTypePassport,
		NationalID:       "123456789",
		IssueDate:        "2020-01-01",
		ExpiryDate:       "2030-01-01",
		AddressLine1:     "123 Main St, #456",
		AddressLine2:     "Building A & B",
		City:             "New York",
		State:            "NY",
		ZipCode:          "10001",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	}

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Logf("SubmitKyc error for special characters in address: %v", err)

		return
	}

	t.Logf("Response for special characters: success=%v, message=%s", resp.Success, resp.Message)
}
