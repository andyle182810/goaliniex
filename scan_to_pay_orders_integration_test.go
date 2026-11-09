package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex/v2"
)

func TestIntegration_CreateScanToPayOrder_ValidUSDTVND(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.ScanToPayOrderRequest{
		OrderID:       getTestOrderID(t, "test-s2p"),
		FiatCurrency:  goaliniex.FiatCurrencyVND,
		FiatAmount:    100000,
		BankCode:      getTestBankCode(t),
		AccountNumber: getTestBankAccountNumber(t),
		Description:   "payment",
	}

	resp, err := client.CreateScanToPayOrder(ctx, email, req)
	if err != nil {
		t.Fatalf("CreateScanToPayOrder returned error: %v", err)
	}

	t.Logf("CreateScanToPayOrder: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("orderId=%s status=%s fiatAmount=%.2f", resp.Data.OrderID, resp.Data.Status, resp.Data.FiatAmount)
	}
}

func TestIntegration_CreateScanToPayOrder_DuplicateOrderID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	orderID := getTestOrderID(t, "test-s2p-dup")

	req := &goaliniex.ScanToPayOrderRequest{
		OrderID:       orderID,
		FiatCurrency:  goaliniex.FiatCurrencyVND,
		FiatAmount:    100000,
		BankCode:      getTestBankCode(t),
		AccountNumber: getTestBankAccountNumber(t),
		Description:   "payment",
	}

	resp1, err := client.CreateScanToPayOrder(ctx, email, req)
	if err != nil {
		t.Fatalf("first CreateScanToPayOrder returned error: %v", err)
	}

	t.Logf("first order: code=%s", resp1.Code)

	resp2, err := client.CreateScanToPayOrder(ctx, email, req)
	if err != nil {
		t.Fatalf("second CreateScanToPayOrder returned error: %v", err)
	}

	if resp2.IsSuccess() {
		t.Error("expected the duplicate orderId to be rejected")
	} else {
		t.Logf("API correctly rejected duplicate orderId: code=%s message=%s", resp2.Code, resp2.Message)
	}
}

func TestIntegration_GetScanToPayOrder(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	orderID := getTestOrderID(t, "test-s2p-details")

	createResp, err := client.CreateScanToPayOrder(ctx, email, &goaliniex.ScanToPayOrderRequest{
		OrderID:       orderID,
		FiatCurrency:  goaliniex.FiatCurrencyVND,
		FiatAmount:    100000,
		BankCode:      getTestBankCode(t),
		AccountNumber: getTestBankAccountNumber(t),
		Description:   "payment",
	})
	if err != nil {
		t.Fatalf("CreateScanToPayOrder returned error: %v", err)
	}

	if !createResp.IsSuccess() {
		t.Skipf("order creation was not successful, skipping lookup: %s", createResp.Message)
	}

	getResp, err := client.GetScanToPayOrder(ctx, email, orderID)
	if err != nil {
		t.Fatalf("GetScanToPayOrder returned error: %v", err)
	}

	if !getResp.IsSuccess() {
		t.Fatalf("expected successful lookup, got code=%s message=%s", getResp.Code, getResp.Message)
	}

	if getResp.Data.OrderID != orderID {
		t.Errorf("OrderID = %q, want %q", getResp.Data.OrderID, orderID)
	}
}

func TestIntegration_CreateScanToPayOrder_InvalidBankCode(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateScanToPayOrder(ctx, email, &goaliniex.ScanToPayOrderRequest{
		OrderID:       getTestOrderID(t, "test-s2p-badbank"),
		FiatCurrency:  goaliniex.FiatCurrencyVND,
		FiatAmount:    100000,
		BankCode:      "INVALID_BANK",
		AccountNumber: getTestBankAccountNumber(t),
		Description:   "payment",
	})
	if err != nil {
		t.Logf("CreateScanToPayOrder returned error (may be expected): %v", err)

		return
	}

	if resp.IsSuccess() {
		t.Log("warning: API accepted an invalid bank code")
	} else {
		t.Logf("API correctly rejected invalid bank code: code=%s message=%s", resp.Code, resp.Message)
	}
}
