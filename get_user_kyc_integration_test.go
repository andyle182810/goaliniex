package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_GetUserKyc_ValidUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: getTestEmail(t),
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
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

	t.Logf("User KYC information retrieved successfully")
	t.Logf("  Name: %s %s", resp.Data.FirstName, resp.Data.LastName)
	t.Logf("  Date of Birth: %s", resp.Data.DateOfBirth)
	t.Logf("  Gender: %s", resp.Data.Gender)
	t.Logf("  Nationality: %s", resp.Data.Nationality)
	t.Logf("  ID Type: %s", resp.Data.IDType)
	t.Logf("  National ID: %s", resp.Data.NationalID)
	t.Logf("  KYC Status: %s", resp.Data.KycStatus)

	if resp.Data.RejectReason != "" {
		t.Logf("  Reject Reason: %s", resp.Data.RejectReason)
	}
}

func TestIntegration_GetUserKyc_KycStatuses(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: getTestEmail(t),
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Fatalf("GetUserKyc returned error: %v", err)
	}

	if !resp.Success {
		t.Logf("Request failed with error code %d: %s", resp.ErrorCode, resp.Message)

		return
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	validStatuses := []string{"NONE", "PROCESSING", "VERIFIED", "REJECTED"}
	statusValid := false

	for _, status := range validStatuses {
		if resp.Data.KycStatus == status {
			statusValid = true

			break
		}
	}

	if !statusValid {
		t.Errorf("unexpected KYC status: %s", resp.Data.KycStatus)
	}

	t.Logf("KYC status: %s", resp.Data.KycStatus)
}

func TestIntegration_GetUserKyc_NonExistentUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "nonexistent_user_12345@example.com",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Logf("GetUserKyc returned error (may be expected): %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API returned success for non-existent user")
	} else {
		t.Logf("API correctly rejected non-existent user: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_GetUserKyc_InvalidEmail(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: "invalid-email-format",
	}

	resp, err := client.GetUserKyc(ctx, req)
	if err != nil {
		t.Logf("GetUserKyc returned error (may be expected): %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted invalid email format")
	} else {
		t.Logf("API correctly rejected invalid email: %s (error code: %d)",
			resp.Message, resp.ErrorCode)
	}
}

func TestIntegration_GetUserKyc_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.GetUserKycRequest{
		UserEmail: getTestEmail(t),
	}

	_, err := client.GetUserKyc(ctx, req)
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

func TestIntegration_GetUserKyc_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.GetUserKycRequest{
		UserEmail: getTestEmail(t),
	}

	cancel()

	_, err := client.GetUserKyc(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}
