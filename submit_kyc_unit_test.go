package goaliniex_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestSubmitKyc_Success(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"id": 123456789,
			"kycStatus": "pending",
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

	resp, err := client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
		UserEmail:        "test@example.com",
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
	})
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}

	if resp.Data.ID != 123456789 {
		t.Errorf("expected id=123456789, got %d", resp.Data.ID)
	}

	if resp.Data.KycStatus != "pending" {
		t.Errorf("expected kycStatus=pending, got %s", resp.Data.KycStatus)
	}

	if resp.Data.Signature != "abc123signature" {
		t.Errorf("expected signature=abc123signature, got %s", resp.Data.Signature)
	}
}

func TestSubmitKyc_WithIDCard(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"id": 987654321,
			"kycStatus": "pending",
			"signature": "xyz789signature"
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

	resp, err := client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
		UserEmail:        "jane@example.com",
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
		AddressLine2:     "123 Test St",
		City:             "Ho Chi Minh",
		State:            "HCM",
		ZipCode:          "70000",
		FrontIDImage:     "base64front",
		BackIDImage:      "base64back",
		HoldIDImage:      "base64hold",
		PhoneNumber:      "1234567890",
		PhoneCountryCode: "1",
	})
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}

	if resp.Data.ID != 987654321 {
		t.Errorf("expected id=987654321, got %d", resp.Data.ID)
	}
}

func TestSubmitKyc_APIError(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": false,
		"message": "User already has KYC submitted",
		"data": null,
		"errorCode": 400
	}`

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, responseBody), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
		UserEmail:        "existing@example.com",
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
	})
	if err != nil {
		t.Fatalf("SubmitKyc returned error: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for duplicate KYC submission")
	}

	if resp.ErrorCode != 400 {
		t.Errorf("expected errorCode=400, got %d", resp.ErrorCode)
	}
}

func TestSubmitKyc_HTTPError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		//nolint:bodyclose // Response body closed by client
		response: mockResponse(http.StatusInternalServerError, "Internal Server Error"),
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
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
	})
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}

	if !errors.Is(err, goaliniex.ErrUnexpectedStatus) {
		t.Errorf("expected ErrUnexpectedStatus, got %v", err)
	}
}

func TestSubmitKyc_NetworkError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: nil,
		err:      errConnectionRefused,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
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
	})
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}

	if !errors.Is(err, goaliniex.ErrHTTPFailure) {
		t.Errorf("expected ErrHTTPFailure, got %v", err)
	}
}

func TestSubmitKyc_ContextCancellation(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: nil,
		err:      context.Canceled,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.SubmitKyc(ctx, &goaliniex.SubmitKycRequest{
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
	})
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestSubmitKyc_ContextTimeout(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.SubmitKyc(ctx, &goaliniex.SubmitKycRequest{
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
	})
	if err == nil {
		t.Fatal("expected error for timeout, got nil")
	}
}

func TestSubmitKyc_InvalidJSON(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, "invalid json response"), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
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
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSubmitKyc_HTTPStatusCodes(t *testing.T) {
	t.Parallel()

	errorCodes := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}

	for _, statusCode := range errorCodes {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			t.Parallel()

			client, err := newTestClientWithMock(&mockHTTPClient{
				//nolint:bodyclose // Response body closed by client
				response: mockResponse(statusCode, `{"error": "test error"}`),
				err:      nil,
			})
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
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
			})
			if err == nil {
				t.Fatalf("expected error for HTTP %d, got nil", statusCode)
			}

			if !errors.Is(err, goaliniex.ErrUnexpectedStatus) {
				t.Errorf("expected ErrUnexpectedStatus for HTTP %d, got %v", statusCode, err)
			}
		})
	}
}

func TestSubmitKyc_GenderValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		gender goaliniex.Gender
	}{
		{"male", goaliniex.GenderMale},
		{"female", goaliniex.GenderFemale},
	}

	for _, test := range testCases { //nolint:dupl
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			responseBody := `{
				"success": true,
				"message": "Success",
				"data": {
					"id": 123456789,
					"kycStatus": "pending",
					"signature": "testsig"
				},
				"errorCode": 0
			}`

			client, err := newTestClientWithMock(&mockHTTPClient{
				//nolint:bodyclose // Response body closed by client
				response: mockResponse(http.StatusOK, responseBody),
				err:      nil,
			})
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			resp, err := client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
				UserEmail:        "test@example.com",
				FirstName:        "Test",
				LastName:         "User",
				DateOfBirth:      "1990-01-01",
				Gender:           test.gender,
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
			})
			if err != nil {
				t.Fatalf("SubmitKyc returned error: %v", err)
			}

			if !resp.Success {
				t.Errorf("expected success=true for gender %s", test.name)
			}
		})
	}
}

func TestSubmitKyc_DocumentTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		docType  goaliniex.IDType
		expected string
	}{
		{"passport", goaliniex.IDTypePassport, "PASSPORT"},
		{"id_card", goaliniex.IDTypeIDCard, "ID_CARD"},
	}

	for _, test := range testCases { //nolint:dupl
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			responseBody := `{
				"success": true,
				"message": "Success",
				"data": {
					"id": 123456789,
					"kycStatus": "pending",
					"signature": "testsig"
				},
				"errorCode": 0
			}`

			client, err := newTestClientWithMock(&mockHTTPClient{
				//nolint:bodyclose // Response body closed by client
				response: mockResponse(http.StatusOK, responseBody),
				err:      nil,
			})
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			resp, err := client.SubmitKyc(context.Background(), &goaliniex.SubmitKycRequest{
				UserEmail:        "test@example.com",
				FirstName:        "Test",
				LastName:         "User",
				DateOfBirth:      "1990-01-01",
				Gender:           goaliniex.GenderMale,
				Nationality:      "US",
				DocumentType:     test.docType,
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
			})
			if err != nil {
				t.Fatalf("SubmitKyc returned error: %v", err)
			}

			if !resp.Success {
				t.Errorf("expected success=true for document type %s", test.name)
			}
		})
	}
}
