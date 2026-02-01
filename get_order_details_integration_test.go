package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_GetOrderDetails_ValidOrder(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	externalOrderID := "order_a0150bb0-2399-40f8-b4c3-bf3765864c73"
	detailsReq := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: externalOrderID,
	}

	resp, err := client.GetOrderDetails(ctx, detailsReq)
	if err != nil {
		t.Fatalf("GetOrderDetails returned error: %v", err)
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

	t.Logf("Order details retrieved successfully")
	t.Logf("  External Order ID: %s", resp.Data.ExternalOrderID)
	t.Logf("  Type: %s", resp.Data.Type)
	t.Logf("  Status: %s", resp.Data.Status)
	t.Logf("  Fiat Amount: %.2f", resp.Data.FiatAmount)
	t.Logf("  Fiat Currency: %s", resp.Data.FiatCurrency)
	t.Logf("  Paid Amount: %.2f", resp.Data.PaidAmount)
	t.Logf("  Created At: %s", resp.Data.CreatedAt)
	t.Logf("  Expires At: %s", resp.Data.ExpiresAt)
}

func TestIntegration_GetOrderDetails_NonExistentOrder(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: "non-existent-order-" + time.Now().Format("20060102150405"),
	}

	resp, err := client.GetOrderDetails(ctx, req)
	if err != nil {
		t.Logf("GetOrderDetails returned error (may be expected): %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API returned success for non-existent order")
	} else {
		t.Logf("API correctly rejected non-existent order: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_GetOrderDetails_EmptyExternalOrderID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: "",
	}

	resp, err := client.GetOrderDetails(ctx, req)
	if err != nil {
		t.Logf("GetOrderDetails returned error for empty ID: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted empty external order ID")
	} else {
		t.Logf("API correctly rejected empty external order ID: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_GetOrderDetails_ResponseDataFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	externalOrderID := "test-fields-detail-" + time.Now().Format("20060102150405")
	createReq := &goaliniex.CreateOrderRequest{
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

	createResp, err := client.CreateOrder(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !createResp.Success {
		t.Skipf("CreateOrder was not successful, skipping: %s", createResp.Message)
	}

	detailsReq := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: externalOrderID,
	}

	resp, err := client.GetOrderDetails(ctx, detailsReq)
	if err != nil {
		t.Fatalf("GetOrderDetails returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful, skipping field validation: %s", resp.Message)
	}

	data := resp.Data

	t.Logf("Order: ExternalOrderID=%q Type=%q Status=%q", data.ExternalOrderID, data.Type, data.Status)
	t.Logf("Amounts: Fiat=%.2f %s Paid=%.2f", data.FiatAmount, data.FiatCurrency, data.PaidAmount)
	t.Logf("Dates: Created=%q Expires=%q", data.CreatedAt, data.ExpiresAt)

	if data.ExternalOrderID != externalOrderID {
		t.Errorf("ExternalOrderID mismatch: got %q, want %q", data.ExternalOrderID, externalOrderID)
	}

	if data.Status == "" {
		t.Error("Status should not be empty")
	}

	if data.CreatedAt == "" {
		t.Error("CreatedAt should not be empty")
	}
}

func TestIntegration_GetOrderDetails_TokenTransferFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	externalOrderID := "test-token-transfer-" + time.Now().Format("20060102150405")
	createReq := &goaliniex.CreateOrderRequest{
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

	createResp, err := client.CreateOrder(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !createResp.Success {
		t.Skipf("CreateOrder was not successful, skipping: %s", createResp.Message)
	}

	detailsReq := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: externalOrderID,
	}

	resp, err := client.GetOrderDetails(ctx, detailsReq)
	if err != nil {
		t.Fatalf("GetOrderDetails returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful: %s", resp.Message)
	}

	token := resp.Data.TokenTransfer

	t.Logf("Token Transfer Details:")
	t.Logf("  Currency: %s", token.Currency)
	t.Logf("  Network: %s", token.Network)
	t.Logf("  Price: %.6f", token.Price)
	t.Logf("  Amount: %.6f", token.Amount)
	t.Logf("  Wallet Address: %s", token.WalletAddress)

	if token.Currency == "" {
		t.Log("Note: TokenTransfer.Currency is empty (may be populated later)")
	}
}

func TestIntegration_GetOrderDetails_BankTransferFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	externalOrderID := "test-bank-transfer-" + time.Now().Format("20060102150405")
	createReq := &goaliniex.CreateOrderRequest{
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

	createResp, err := client.CreateOrder(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if !createResp.Success {
		t.Skipf("CreateOrder was not successful, skipping: %s", createResp.Message)
	}

	detailsReq := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: externalOrderID,
	}

	resp, err := client.GetOrderDetails(ctx, detailsReq)
	if err != nil {
		t.Fatalf("GetOrderDetails returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful: %s", resp.Message)
	}

	bank := resp.Data.BankTransfer

	t.Logf("Bank Transfer Details:")
	t.Logf("  Bank Code: %s", bank.BankCode)
	t.Logf("  Bank Name: %s", bank.BankName)
	t.Logf("  Account Number: %s", bank.BankAccountNumber)
	t.Logf("  Account Name: %s", bank.BankAccountName)
	t.Logf("  Content: %s", bank.Content)
	t.Logf("  QR Code URL: %s", bank.QRCodeURL)
}

func TestIntegration_GetOrderDetails_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: "test-timeout-" + time.Now().Format("20060102150405"),
	}

	_, err := client.GetOrderDetails(ctx, req)
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

func TestIntegration_GetOrderDetails_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.GetOrderDetailsRequest{
		ExternalOrderID: "test-cancel-" + time.Now().Format("20060102150405"),
	}

	cancel()

	_, err := client.GetOrderDetails(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}

func TestIntegration_GetOrderDetails_DifferentCryptoCurrencies(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	webhookSecret := getWebhookSecretKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	t.Cleanup(cancel)

	currencies := []goaliniex.Currency{
		goaliniex.CurrencyUSDT,
		goaliniex.CurrencyETH,
		goaliniex.CurrencyBTC,
	}

	for _, currency := range currencies {
		t.Run(string(currency), func(t *testing.T) {
			t.Parallel()

			externalOrderID := "test-" + string(currency) + "-details-" + time.Now().Format("20060102150405")
			createReq := &goaliniex.CreateOrderRequest{
				Currency:          currency,
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

			createResp, err := client.CreateOrder(ctx, createReq)
			if err != nil {
				t.Logf("CreateOrder error for %s: %v", currency, err)

				return
			}

			if !createResp.Success {
				t.Logf("CreateOrder failed for %s: %s", currency, createResp.Message)

				return
			}

			detailsReq := &goaliniex.GetOrderDetailsRequest{
				ExternalOrderID: externalOrderID,
			}

			resp, err := client.GetOrderDetails(ctx, detailsReq)
			if err != nil {
				t.Logf("GetOrderDetails error for %s: %v", currency, err)

				return
			}

			t.Logf("Response for %s: success=%v, errorCode=%d, message=%s",
				currency, resp.Success, resp.ErrorCode, resp.Message)

			if resp.Success && resp.Data != nil {
				t.Logf("  Status: %s, TokenCurrency: %s",
					resp.Data.Status, resp.Data.TokenTransfer.Currency)
			}
		})
	}
}
