package goaliniex_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestGetWalletBalance_Success(t *testing.T) {
	t.Parallel()

	successResponse := `{
		"success": true,
		"message": "Success",
		"data": {
			"balance": 1234.56,
			"currency": "USDT",
			"signature": "mock-signature"
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

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyUSDT,
	}

	resp, err := client.GetWalletBalance(ctx, req)
	if err != nil {
		t.Fatalf("GetWalletBalance returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.Balance != 1234.56 {
		t.Errorf("expected balance=1234.56, got %f", resp.Data.Balance)
	}

	if resp.Data.Currency != goaliniex.CurrencyUSDT {
		t.Errorf("expected currency=USDT, got %s", resp.Data.Currency)
	}

	if resp.Data.Signature != "mock-signature" {
		t.Errorf("expected signature=mock-signature, got %s", resp.Data.Signature)
	}
}

func TestGetWalletBalance_APIError(t *testing.T) {
	t.Parallel()

	errorResponse := `{
		"success": false,
		"message": "Invalid currency",
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

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.Currency("INVALID"),
	}

	resp, err := client.GetWalletBalance(ctx, req)
	if err != nil {
		t.Fatalf("GetWalletBalance returned error: %v", err)
	}

	if resp.Success {
		t.Error("expected success=false, got true")
	}

	if resp.ErrorCode != 1001 {
		t.Errorf("expected errorCode=1001, got %d", resp.ErrorCode)
	}

	if resp.Message != "Invalid currency" {
		t.Errorf("expected message='Invalid currency', got %s", resp.Message)
	}
}

func TestGetWalletBalance_HTTPError(t *testing.T) {
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

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyUSDT,
	}

	_, err = client.GetWalletBalance(ctx, req)
	if err == nil {
		t.Fatal("expected error for HTTP 400, got nil")
	}

	t.Logf("Got expected error: %v", err)
}

func TestGetWalletBalance_ZeroBalance(t *testing.T) {
	t.Parallel()

	zeroBalanceResponse := `{
		"success": true,
		"message": "Success",
		"data": {
			"balance": 0,
			"currency": "BTC",
			"signature": "sig123"
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, zeroBalanceResponse), //nolint:bodyclose
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyBTC,
	}

	resp, err := client.GetWalletBalance(ctx, req)
	if err != nil {
		t.Fatalf("GetWalletBalance returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false: %s", resp.Message)
	}

	if resp.Data.Balance != 0 {
		t.Errorf("expected balance=0, got %f", resp.Data.Balance)
	}
}
