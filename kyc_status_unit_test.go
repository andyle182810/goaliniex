package goaliniex_test

import (
	"context"
	"testing"

	"github.com/andyle182810/goaliniex/v2"
)

func TestGetKycStatus_Success(t *testing.T) {
	t.Parallel()

	client := newTestClientForAuthenticatedCall(
		`{"code":"SUCCESS","message":"ok","data":{"status":"VERIFIED"}}`,
		200,
	)

	resp, err := client.GetKycStatus(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("GetKycStatus: %v", err)
	}

	if !resp.IsSuccess() {
		t.Fatalf("expected success, got code=%s message=%s", resp.Code, resp.Message)
	}

	if resp.Data.Status != goaliniex.KycStatusVerified {
		t.Errorf("Status = %q, want %q", resp.Data.Status, goaliniex.KycStatusVerified)
	}
}

func TestGetKycStatus_BusinessError(t *testing.T) {
	t.Parallel()

	client := newTestClientForAuthenticatedCall(
		`{"code":"PARTNER_ACCOUNT_LOCKED","message":"locked","data":null}`,
		403,
	)

	resp, err := client.GetKycStatus(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("GetKycStatus: %v", err)
	}

	if resp.IsSuccess() {
		t.Fatal("expected a business error, got success")
	}

	if resp.Code != goaliniex.CodePartnerAccountLocked {
		t.Errorf("Code = %q, want %q", resp.Code, goaliniex.CodePartnerAccountLocked)
	}
}
