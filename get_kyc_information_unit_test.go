package goaliniex_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

var errConnectionRefused = errors.New("connection refused")

func TestGetKycInformation_Success(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "John",
			"lastName": "Doe",
			"dateOfBirth": "1990-01-01",
			"gender": "male",
			"nationality": "US",
			"idType": "passport",
			"nationalId": "123456789",
			"issueDate": "2020-01-01",
			"expiryDate": "2030-01-01",
			"address": "123 Main St",
			"frontIdImage": "base64image",
			"backIdImage": "base64image",
			"holdIdImage": "base64image",
			"phoneNumber": "1234567890",
			"phoneCountryCode": "1",
			"kycStatus": "approved",
			"rejectReason": ""
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

	resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}

	if resp.Data.FirstName != "John" {
		t.Errorf("expected firstName=John, got %s", resp.Data.FirstName)
	}

	if resp.Data.LastName != "Doe" {
		t.Errorf("expected lastName=Doe, got %s", resp.Data.LastName)
	}

	if resp.Data.KycStatus != "approved" {
		t.Errorf("expected kycStatus=approved, got %s", resp.Data.KycStatus)
	}
}

func TestGetKycInformation_PendingStatus(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "Jane",
			"lastName": "Smith",
			"dateOfBirth": "1985-05-15",
			"gender": "female",
			"nationality": "VN",
			"idType": "id_card",
			"nationalId": "987654321",
			"issueDate": "2019-06-01",
			"expiryDate": "2029-06-01",
			"address": "456 Oak Ave",
			"frontIdImage": "",
			"backIdImage": "",
			"holdIdImage": "",
			"phoneNumber": "0987654321",
			"phoneCountryCode": "84",
			"kycStatus": "pending",
			"rejectReason": ""
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

	resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "jane@example.com",
	})
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if resp.Data.KycStatus != "pending" {
		t.Errorf("expected kycStatus=pending, got %s", resp.Data.KycStatus)
	}
}

func TestGetKycInformation_RejectedStatus(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "Bob",
			"lastName": "Wilson",
			"dateOfBirth": "1992-12-25",
			"gender": "male",
			"nationality": "VN",
			"idType": "passport",
			"nationalId": "ABC123456",
			"issueDate": "2018-01-01",
			"expiryDate": "2028-01-01",
			"address": "789 Pine St",
			"frontIdImage": "",
			"backIdImage": "",
			"holdIdImage": "",
			"phoneNumber": "0123456789",
			"phoneCountryCode": "84",
			"kycStatus": "rejected",
			"rejectReason": "Document is blurry"
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

	resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "bob@example.com",
	})
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if resp.Data.KycStatus != "rejected" {
		t.Errorf("expected kycStatus=rejected, got %s", resp.Data.KycStatus)
	}

	if resp.Data.RejectReason != "Document is blurry" {
		t.Errorf("expected rejectReason='Document is blurry', got %s", resp.Data.RejectReason)
	}
}

func TestGetKycInformation_APIError(t *testing.T) {
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

	resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "nonexistent@example.com",
	})
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false for non-existent user")
	}

	if resp.ErrorCode != 404 {
		t.Errorf("expected errorCode=404, got %d", resp.ErrorCode)
	}
}

func TestGetKycInformation_HTTPError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		//nolint:bodyclose // Response body closed by client
		response: mockResponse(http.StatusInternalServerError, "Internal Server Error"),
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}

	if !errors.Is(err, goaliniex.ErrUnexpectedStatus) {
		t.Errorf("expected ErrUnexpectedStatus, got %v", err)
	}
}

func TestGetKycInformation_NetworkError(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: nil,
		err:      errConnectionRefused,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}

	if !errors.Is(err, goaliniex.ErrHTTPFailure) {
		t.Errorf("expected ErrHTTPFailure, got %v", err)
	}
}

func TestGetKycInformation_ContextCancellation(t *testing.T) {
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

	_, err = client.GetKycInformation(ctx, &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestGetKycInformation_ContextTimeout(t *testing.T) {
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

	_, err = client.GetKycInformation(ctx, &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error for timeout, got nil")
	}
}

func TestGetKycInformation_InvalidJSON(t *testing.T) {
	t.Parallel()

	client, err := newTestClientWithMock(&mockHTTPClient{
		response: mockResponse(http.StatusOK, "invalid json response"), //nolint:bodyclose // Response body closed by client
		err:      nil,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGetKycInformation_EmptyResponse(t *testing.T) {
	t.Parallel()

	responseBody := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "",
			"lastName": "",
			"dateOfBirth": "",
			"gender": "",
			"nationality": "",
			"idType": "",
			"nationalId": "",
			"issueDate": "",
			"expiryDate": "",
			"address": "",
			"frontIdImage": "",
			"backIdImage": "",
			"holdIdImage": "",
			"phoneNumber": "",
			"phoneCountryCode": "",
			"kycStatus": "not_submitted",
			"rejectReason": ""
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

	resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
		UserEmail: "new_user@example.com",
	})
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}

	if resp.Data.KycStatus != "not_submitted" {
		t.Errorf("expected kycStatus=not_submitted, got %s", resp.Data.KycStatus)
	}
}

func TestGetKycInformation_AllKycStatuses(t *testing.T) {
	t.Parallel()

	testCases := []goaliniex.KycStatus{
		goaliniex.KycStatusNone,
		goaliniex.KycStatusProcessing,
		goaliniex.KycStatusVerified,
		goaliniex.KycStatusRejected,
	}

	for _, status := range testCases {
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()

			responseBody := `{
				"success": true,
				"message": "Success",
				"data": {
					"firstName": "Test",
					"lastName": "User",
					"dateOfBirth": "",
					"gender": "",
					"nationality": "",
					"idType": "",
					"nationalId": "",
					"issueDate": "",
					"expiryDate": "",
					"address": "",
					"frontIdImage": "",
					"backIdImage": "",
					"holdIdImage": "",
					"phoneNumber": "",
					"phoneCountryCode": "",
					"kycStatus": "` + string(status) + `",
					"rejectReason": ""
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

			resp, err := client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
				UserEmail: "test@example.com",
			})
			if err != nil {
				t.Fatalf("GetKycInformation returned error: %v", err)
			}

			if resp.Data.KycStatus != status {
				t.Errorf("expected kycStatus=%s, got %s", status, resp.Data.KycStatus)
			}
		})
	}
}

func TestGetKycInformation_HTTPStatusCodes(t *testing.T) {
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

			_, err = client.GetKycInformation(context.Background(), &goaliniex.KycInformationRequest{
				UserEmail: "test@example.com",
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
