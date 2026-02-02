package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_GetWalletBalance_USDT(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyUSDT,
	}

	resp, err := client.GetWalletBalance(ctx, req)
	if err != nil {
		t.Fatalf("GetWalletBalance returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if !resp.Success {
		t.Logf("Request failed with error code %d: %s", resp.ErrorCode, resp.Message)

		return
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	t.Logf("Wallet balance retrieved successfully")
	t.Logf("  Balance: %.6f", resp.Data.Balance)
	t.Logf("  Currency: %s", resp.Data.Currency)
	t.Logf("  Signature: %s", resp.Data.Signature)
}

func TestIntegration_GetWalletBalance_AllCurrencies(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	currencies := []goaliniex.Currency{
		goaliniex.CurrencyUSDT,
		goaliniex.CurrencyETH,
		goaliniex.CurrencyBTC,
	}

	for _, currency := range currencies {
		t.Run(string(currency), func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.GetWalletBalanceRequest{
				Currency: currency,
			}

			resp, err := client.GetWalletBalance(ctx, req)
			if err != nil {
				t.Logf("GetWalletBalance error for %s: %v", currency, err)

				return
			}

			if !resp.Success {
				t.Logf("Request failed for %s: error code %d: %s",
					currency, resp.ErrorCode, resp.Message)

				return
			}

			if resp.Data != nil {
				t.Logf("%s balance: %.6f", currency, resp.Data.Balance)
			}
		})
	}
}

func TestIntegration_GetWalletBalance_InvalidCurrency(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.Currency("INVALID"),
	}

	resp, err := client.GetWalletBalance(ctx, req)
	if err != nil {
		t.Logf("GetWalletBalance returned error (may be expected): %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted invalid currency")
	} else {
		t.Logf("API correctly rejected invalid currency: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_GetWalletBalance_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyUSDT,
	}

	_, err := client.GetWalletBalance(ctx, req)
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

func TestIntegration_GetWalletBalance_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.GetWalletBalanceRequest{
		Currency: goaliniex.CurrencyUSDT,
	}

	cancel()

	_, err := client.GetWalletBalance(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}
