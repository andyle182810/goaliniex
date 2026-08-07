package goaliniex_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/andyle182810/goaliniex"
)

func TestSubmitKycSupplement_Success(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"id": 123456789,
			"kycStatus": "VERIFIED",
			"kycSupplementStatus": "SUBMITTED",
			"signature": "abc123signature"
		},
		"errorCode": 0
	}`

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, responseBody), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.SubmitKycSupplement(context.Background(), &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        "test@example.com",
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
		KycReport:        "data:application/pdf;base64,report",
		PhoneNumber:      "",
		PhoneCountryCode: "",
	})
	if err != nil {
		t.Fatalf("SubmitKycSupplement returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}

	if resp.Data.ID != 123456789 {
		t.Errorf("expected id=123456789, got %d", resp.Data.ID)
	}

	if resp.Data.KycStatus != "VERIFIED" {
		t.Errorf("expected kycStatus=VERIFIED, got %s", resp.Data.KycStatus)
	}

	if resp.Data.KycSupplementStatus != "SUBMITTED" {
		t.Errorf("expected kycSupplementStatus=SUBMITTED, got %s", resp.Data.KycSupplementStatus)
	}
}

func TestSubmitKycSupplement_BackfillImagesOnly(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"id": 123456789,
			"kycStatus": "VERIFIED",
			"kycSupplementStatus": "SUBMITTED",
			"signature": "abc123signature"
		},
		"errorCode": 0
	}`

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, responseBody), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.SubmitKycSupplement(context.Background(), &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        "test@example.com",
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
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		KycReport:        "",
		PhoneNumber:      "",
		PhoneCountryCode: "",
	})
	if err != nil {
		t.Fatalf("SubmitKycSupplement returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}
}

func TestSubmitKycSupplement_APIError(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": false,
		"message": "User not found",
		"data": null,
		"errorCode": 404
	}`

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, responseBody), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.SubmitKycSupplement(context.Background(), &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        "nonexistent@example.com",
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
		KycReport:        "data:application/pdf;base64,report",
		PhoneNumber:      "",
		PhoneCountryCode: "",
	})
	if err != nil {
		t.Fatalf("SubmitKycSupplement returned error: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for nonexistent user")
	}

	if resp.ErrorCode != 404 {
		t.Errorf("expected errorCode=404, got %d", resp.ErrorCode)
	}
}

func TestSubmitKycSupplement_HTTPError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		//nolint:bodyclose // Response body closed by client
		response: mockResponse(http.StatusInternalServerError, "Internal Server Error"),
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.SubmitKycSupplement(context.Background(), &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        "test@example.com",
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
		KycReport:        "",
		PhoneNumber:      "",
		PhoneCountryCode: "",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}

	if !errors.Is(err, goaliniex.ErrUnexpectedStatus) {
		t.Errorf("expected ErrUnexpectedStatus, got %v", err)
	}
}

func TestSubmitKycSupplement_NetworkError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: nil,
		err:      errConnectionRefused,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.SubmitKycSupplement(context.Background(), &goaliniex.SubmitKycSupplementRequest{
		UserEmail:        "test@example.com",
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
		KycReport:        "",
		PhoneNumber:      "",
		PhoneCountryCode: "",
	})
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}

	if !errors.Is(err, goaliniex.ErrHTTPFailure) {
		t.Errorf("expected ErrHTTPFailure, got %v", err)
	}
}
