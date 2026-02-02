package goaliniex_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func getWebhookSecretKey(t *testing.T) string {
	t.Helper()

	webhookSecret := strings.TrimSpace(os.Getenv("ALIX_WEBHOOK_SECRET_KEY"))
	if webhookSecret == "" {
		t.Skip("skipping test: ALIX_WEBHOOK_SECRET_KEY not set")
	}

	return webhookSecret
}

func TestIntegration_CreateOrder_ValidUSDTVND(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          getTestBankCode(t),
		BankAccountNumber: getTestBankAccountNumber(t),
		ExternalOrderID:   "test-order-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
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

	t.Logf("Order created successfully")
	t.Logf("  External Order ID: %s", resp.Data.ExternalOrderID)
	t.Logf("  Type: %s", resp.Data.Type)
	t.Logf("  Status: %s", resp.Data.Status)
	t.Logf("  Fiat Amount: %.2f", resp.Data.FiatAmount)
	t.Logf("  Token: %s, Price: %.2f, Amount: %.4f",
		resp.Data.TokenTransfer.Currency, resp.Data.TokenTransfer.Price, resp.Data.TokenTransfer.Amount)
	t.Logf("  Bank: %s (%s)", resp.Data.BankTransfer.BankName, resp.Data.BankTransfer.BankCode)
	t.Logf("  Bank Account: %s (%s)", resp.Data.BankTransfer.BankAccountNumber, resp.Data.BankTransfer.BankAccountName)
	t.Logf("  Fees: System=%.2f, Processing=%.2f", resp.Data.Fees.SystemFee, resp.Data.Fees.ProcessingFee)
	t.Logf("  Created At: %s, Expires At: %s", resp.Data.CreatedAt, resp.Data.ExpiresAt)
}

func TestIntegration_CreateOrder_WithKYCVerified(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        50000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-kyc-order-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !resp.Success {
		t.Logf("Request failed with error code %d: %s", resp.ErrorCode, resp.Message)

		return
	}

	t.Logf("Order created with KYC verified user")
	t.Logf("  External Order ID: %s", resp.Data.ExternalOrderID)
	t.Logf("  Status: %s", resp.Data.Status)
}

func TestIntegration_CreateOrder_DifferentFiatCurrencies(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	testCases := []struct {
		name         string
		fiatCurrency goaliniex.FiatCurrency
		fiatAmount   float64
		bankCode     string
	}{
		{
			name:         "VND",
			fiatCurrency: goaliniex.FiatCurrencyVND,
			fiatAmount:   100000,
			bankCode:     "970407",
		},
		{
			name:         "PHP",
			fiatCurrency: goaliniex.FiatCurrencyPHP,
			fiatAmount:   1000,
			bankCode:     "TESTBANK",
		},
		{
			name:         "THB",
			fiatCurrency: goaliniex.FiatCurrencyTHB,
			fiatAmount:   1000,
			bankCode:     "TESTBANK",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.CreateOrderRequest{
				Currency:          goaliniex.CurrencyUSDT,
				FiatAmount:        testCase.fiatAmount,
				FiatCurrency:      testCase.fiatCurrency,
				BankCode:          testCase.bankCode,
				BankAccountNumber: "123456789",
				ExternalOrderID:   "test-" + testCase.name + "-" + time.Now().Format("20060102150405"),
				WebhookSecretKey:  webhookSecret,
				UserEmail:         getTestEmail(t),
				UserKYCVerified:   true,
				Content:           "payment",
				ExtendInfo:        nil,
			}

			resp, err := client.CreateOrder(ctx, req)
			if err != nil {
				t.Logf("CreateOrder error for %s: %v", testCase.name, err)

				return
			}

			t.Logf("Response for %s: success=%v, errorCode=%d, message=%s",
				testCase.name, resp.Success, resp.ErrorCode, resp.Message)

			if resp.Success && resp.Data != nil {
				t.Logf("  Order ID: %s, Status: %s", resp.Data.ExternalOrderID, resp.Data.Status)
			}
		})
	}
}

func TestIntegration_CreateOrder_DifferentCryptoCurrencies(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

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

			req := &goaliniex.CreateOrderRequest{
				Currency:          currency,
				FiatAmount:        100000,
				FiatCurrency:      goaliniex.FiatCurrencyVND,
				BankCode:          "970407",
				BankAccountNumber: "888812345678",
				ExternalOrderID:   "test-" + string(currency) + "-" + time.Now().Format("20060102150405"),
				WebhookSecretKey:  webhookSecret,
				UserEmail:         getTestEmail(t),
				UserKYCVerified:   true,
				Content:           "payment",
				ExtendInfo:        nil,
			}

			resp, err := client.CreateOrder(ctx, req)
			if err != nil {
				t.Logf("CreateOrder error for %s: %v", currency, err)

				return
			}

			t.Logf("Response for %s: success=%v, errorCode=%d, message=%s",
				currency, resp.Success, resp.ErrorCode, resp.Message)

			if resp.Success && resp.Data != nil {
				t.Logf("  Order ID: %s, Status: %s", resp.Data.ExternalOrderID, resp.Data.Status)
			}
		})
	}
}

func TestIntegration_CreateOrder_InvalidBankCode(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "INVALID_BANK",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-invalid-bank-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("CreateOrder returned error (may be expected): %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted invalid bank code")
	} else {
		t.Logf("API correctly rejected invalid bank code: %s (error code: %d)", resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_CreateOrder_ZeroAmount(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        0,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-zero-amount-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("CreateOrder returned error for zero amount: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted zero amount order")
	} else {
		t.Logf("API correctly rejected zero amount: %s (error code: %d)", resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_CreateOrder_NegativeAmount(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        -100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-negative-amount-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("CreateOrder returned error for negative amount: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted negative amount order")
	} else {
		t.Logf("API correctly rejected negative amount: %s (error code: %d)", resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_CreateOrder_DuplicateExternalOrderID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	externalOrderID := "test-duplicate-" + time.Now().Format("20060102150405")

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   externalOrderID,
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp1, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("First CreateOrder returned error: %v", err)
	}

	t.Logf("First order: success=%v, errorCode=%d", resp1.Success, resp1.ErrorCode)

	resp2, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("Second CreateOrder returned error (may be expected): %v", err)

		return
	}

	if resp2.Success {
		t.Logf("Warning: API accepted duplicate external order ID")
	} else {
		t.Logf("API correctly rejected duplicate external order ID: %s (error code: %d)",
			resp2.Message, resp2.ErrorCode)
	}
}

func TestIntegration_CreateOrder_EmptyExternalOrderID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "",
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("CreateOrder returned error for empty external order ID: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted empty external order ID")
	} else {
		t.Logf("API correctly rejected empty external order ID: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_CreateOrder_InvalidEmail(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-invalid-email-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         "not-an-email",
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Logf("CreateOrder returned error for invalid email: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted invalid email format")
	} else {
		t.Logf("API correctly rejected invalid email: %s (error code: %d)", resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_CreateOrder_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-timeout-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	_, err := client.CreateOrder(ctx, req)
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

func TestIntegration_CreateOrder_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-cancel-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	cancel()

	_, err := client.CreateOrder(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}

func TestIntegration_CreateOrder_ResponseDataFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-fields-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful, skipping field validation: %s", resp.Message)
	}

	data := resp.Data

	t.Logf("Response Data Fields:")
	t.Logf("  ExternalOrderID: %q", data.ExternalOrderID)
	t.Logf("  Type: %q", data.Type)
	t.Logf("  Status: %q", data.Status)
	t.Logf("  FiatAmount: %.2f", data.FiatAmount)
	t.Logf("  TokenTransfer: Currency=%q, Price=%.2f, Amount=%.4f",
		data.TokenTransfer.Currency, data.TokenTransfer.Price, data.TokenTransfer.Amount)
	t.Logf("  BankTransfer: BankCode=%q, BankName=%q, AccountNumber=%q",
		data.BankTransfer.BankCode, data.BankTransfer.BankName, data.BankTransfer.BankAccountNumber)
	t.Logf("  Fees: SystemFee=%.2f, ProcessingFee=%.2f", data.Fees.SystemFee, data.Fees.ProcessingFee)
	t.Logf("  CreatedAt: %q", data.CreatedAt)
	t.Logf("  ExpiresAt: %q", data.ExpiresAt)

	if data.ExternalOrderID != req.ExternalOrderID {
		t.Errorf("ExternalOrderID mismatch: got %q, want %q", data.ExternalOrderID, req.ExternalOrderID)
	}

	if data.Status == "" {
		t.Error("Status should not be empty")
	}

	if data.TokenTransfer.Currency != req.Currency {
		t.Errorf("Currency mismatch: got %q, want %q", data.TokenTransfer.Currency, req.Currency)
	}

	if data.BankTransfer.BankCode == "" {
		t.Error("BankCode should not be empty")
	}
}

func TestIntegration_CreateOrder_LongTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.CreateOrderRequest{
		Currency:          goaliniex.CurrencyUSDT,
		FiatAmount:        100000,
		FiatCurrency:      goaliniex.FiatCurrencyVND,
		BankCode:          "970407",
		BankAccountNumber: "888812345678",
		ExternalOrderID:   "test-long-timeout-" + time.Now().Format("20060102150405"),
		WebhookSecretKey:  webhookSecret,
		UserEmail:         getTestEmail(t),
		UserKYCVerified:   true,
		Content:           "payment",
		ExtendInfo:        nil,
	}

	start := time.Now()
	resp, err := client.CreateOrder(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	t.Logf("Request completed in %v", duration)
	t.Logf("Response: success=%v, errorCode=%d, message=%s",
		resp.Success, resp.ErrorCode, resp.Message)

	if resp.Success && resp.Data != nil {
		t.Logf("  Order ID: %s, Status: %s", resp.Data.ExternalOrderID, resp.Data.Status)
	}
}
