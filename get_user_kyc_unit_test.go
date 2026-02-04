package goaliniex_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestGetUserKyc_Success(t *testing.T) {
	t.Parallel()

	successResponse := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "John",
			"lastName": "Doe",
			"dateOfBirth": "1990-01-15",
			"gender": "Male",
			"nationality": "US",
			"idType": "PASSPORT",
			"nationalId": "A12345678",
			"issueDate": "2020-01-01",
			"expiryDate": "2030-01-01",
			"address": "123 Main St, New York, NY",
			"frontIdImage": "base64encodedimage",
			"backIdImage": "base64encodedimage",
			"holdIdImage": "base64encodedimage",
			"phoneNumber": "1234567890",
			"phoneCountryCode": "+1",
			"kycStatus": "VERIFIED"
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, successResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "john.doe@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.FirstName != "John" {
		t.Errorf("expected firstName=John, got %s", resp.Data.FirstName)
	}

	if resp.Data.LastName != "Doe" {
		t.Errorf("expected lastName=Doe, got %s", resp.Data.LastName)
	}

	if resp.Data.KycStatus != "VERIFIED" {
		t.Errorf("expected kycStatus=VERIFIED, got %s", resp.Data.KycStatus)
	}

	if resp.Data.IDType != "PASSPORT" {
		t.Errorf("expected idType=PASSPORT, got %s", resp.Data.IDType)
	}
}

func TestGetUserKyc_Rejected(t *testing.T) {
	t.Parallel()

	rejectedResponse := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "Jane",
			"lastName": "Smith",
			"dateOfBirth": "1985-06-20",
			"gender": "Female",
			"nationality": "VN",
			"idType": "ID_CARD",
			"nationalId": "123456789012",
			"issueDate": "2018-05-10",
			"expiryDate": "2028-05-10",
			"address": "456 Le Loi, Ho Chi Minh City",
			"frontIdImage": "base64encodedimage",
			"backIdImage": "base64encodedimage",
			"holdIdImage": "base64encodedimage",
			"phoneNumber": "987654321",
			"phoneCountryCode": "+84",
			"kycStatus": "REJECTED",
			"rejectReason": "ID document expired"
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, rejectedResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "jane.smith@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.KycStatus != "REJECTED" {
		t.Errorf("expected kycStatus=REJECTED, got %s", resp.Data.KycStatus)
	}

	if resp.Data.RejectReason != "ID document expired" {
		t.Errorf("expected rejectReason='ID document expired', got %s", resp.Data.RejectReason)
	}
}

func TestGetUserKyc_Processing(t *testing.T) {
	t.Parallel()

	processingResponse := `{
		"success": true,
		"message": "Success",
		"data": {
			"firstName": "Bob",
			"lastName": "Johnson",
			"dateOfBirth": "1992-03-25",
			"gender": "Male",
			"nationality": "US",
			"idType": "ID_CARD",
			"nationalId": "987654321",
			"issueDate": "2022-01-01",
			"expiryDate": "2032-01-01",
			"address": "789 Oak Ave, Los Angeles, CA",
			"frontIdImage": "base64encodedimage",
			"backIdImage": "base64encodedimage",
			"holdIdImage": "base64encodedimage",
			"phoneNumber": "5551234567",
			"phoneCountryCode": "+1",
			"kycStatus": "PROCESSING"
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, processingResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "bob.johnson@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.KycStatus != "PROCESSING" {
		t.Errorf("expected kycStatus=PROCESSING, got %s", resp.Data.KycStatus)
	}
}

func TestGetUserKyc_UserNotFound(t *testing.T) {
	t.Parallel()

	errorResponse := `{
		"success": false,
		"message": "User not found",
		"data": null,
		"errorCode": 1001
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, errorResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "nonexistent@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false, got true")
	}

	if resp.ErrorCode != 1001 {
		t.Errorf("expected errorCode=1001, got %d", resp.ErrorCode)
	}

	if resp.Message != "User not found" {
		t.Errorf("expected message='User not found', got %s", resp.Message)
	}
}

func TestGetUserKyc_HTTPError(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusBadRequest, `{"error": "Bad Request"}`), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "test@example.com",
	}

	_, err = client.GetUserKyc(ctx, req)
	if err == nil {
		t.Fatal("expected error for HTTP 400, got nil")
	}

	t.Logf("Got expected error: %v", err)
}

func TestGetUserKyc_NoKyc(t *testing.T) {
	t.Parallel()

	noKycResponse := `{
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
			"kycStatus": "NONE"
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, noKycResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "newuser@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.KycStatus != "NONE" {
		t.Errorf("expected kycStatus=NONE, got %s", resp.Data.KycStatus)
	}
}
