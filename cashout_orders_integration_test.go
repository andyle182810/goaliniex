package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex/v2"
)

func TestIntegration_CreateCashoutOrder_ValidUSDTVND(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CashoutOrderRequest{
		OrderID:      getTestOrderID(t, "test-cashout"),
		FiatCurrency: goaliniex.FiatCurrencyVND,
		FiatAmount:   100000,
		Recipient: goaliniex.CashoutRecipient{
			BankCode:      getTestBankCode(t),
			AccountNumber: getTestBankAccountNumber(t),
			AccountName:   "Test User",
			Description:   "payment",
		},
	}

	resp, err := client.CreateCashoutOrder(ctx, email, req)
	if err != nil {
		t.Fatalf("CreateCashoutOrder returned error: %v", err)
	}

	t.Logf("CreateCashoutOrder: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("orderId=%s status=%s", resp.Data.OrderID, resp.Data.Status)
	}
}

// AccountName is required for cashout (unlike scan-to-pay); omitting it
// should be rejected by the server with a structured error code.
func TestIntegration_CreateCashoutOrder_MissingAccountName(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateCashoutOrder(ctx, email, &goaliniex.CashoutOrderRequest{
		OrderID:      getTestOrderID(t, "test-cashout-noname"),
		FiatCurrency: goaliniex.FiatCurrencyVND,
		FiatAmount:   100000,
		Recipient: goaliniex.CashoutRecipient{
			BankCode:      getTestBankCode(t),
			AccountNumber: getTestBankAccountNumber(t),
			Description:   "payment",
		},
	})
	if err != nil {
		t.Logf("CreateCashoutOrder returned error (may be expected): %v", err)

		return
	}

	if resp.IsSuccess() {
		t.Log("warning: API accepted a cashout order without accountName")
	} else {
		t.Logf("API correctly rejected missing accountName: code=%s message=%s", resp.Code, resp.Message)
	}
}

func TestIntegration_GetCashoutOrder(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	email := getTestEmail(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	orderID := getTestOrderID(t, "test-cashout-details")

	createResp, err := client.CreateCashoutOrder(ctx, email, &goaliniex.CashoutOrderRequest{
		OrderID:      orderID,
		FiatCurrency: goaliniex.FiatCurrencyVND,
		FiatAmount:   100000,
		Recipient: goaliniex.CashoutRecipient{
			BankCode:      getTestBankCode(t),
			AccountNumber: getTestBankAccountNumber(t),
			AccountName:   "Test User",
			Description:   "payment",
		},
	})
	if err != nil {
		t.Fatalf("CreateCashoutOrder returned error: %v", err)
	}

	if !createResp.IsSuccess() {
		t.Skipf("order creation was not successful, skipping lookup: %s", createResp.Message)
	}

	getResp, err := client.GetCashoutOrder(ctx, email, orderID)
	if err != nil {
		t.Fatalf("GetCashoutOrder returned error: %v", err)
	}

	if !getResp.IsSuccess() {
		t.Fatalf("expected successful lookup, got code=%s message=%s", getResp.Code, getResp.Message)
	}

	if getResp.Data.OrderID != orderID {
		t.Errorf("OrderID = %q, want %q", getResp.Data.OrderID, orderID)
	}
}
