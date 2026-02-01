package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_GetKycInformation_ValidUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: getTestEmail(t),
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if !resp.Success {
		t.Fatalf("request failed: %+v", resp)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	t.Logf("KYC Information retrieved successfully")
	t.Logf("  First Name: %s", resp.Data.FirstName)
	t.Logf("  Last Name: %s", resp.Data.LastName)
	t.Logf("  KYC Status: %s", resp.Data.KycStatus)
}

func TestIntegration_GetKycInformation_NonExistentUser(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: "nonexistent_user_12345@example.com",
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Logf("GetKycInformation returned error (may be expected): %v", err)

		return
	}

	if resp.Success && resp.Data != nil {
		t.Logf("API returned data for non-existent user (may be expected behavior)")
	} else {
		t.Logf("API correctly indicated user not found or returned empty data")
	}
}

func TestIntegration_GetKycInformation_InvalidEmailFormat(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: "invalid-email-format",
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Logf("GetKycInformation returned error for invalid email: %v", err)

		return
	}

	t.Logf("Response success: %v, message: %s", resp.Success, resp.Message)
}

func TestIntegration_GetKycInformation_EmptyEmail(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: "",
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Logf("GetKycInformation returned error for empty email: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted empty email, this might be a validation issue")
	} else {
		t.Logf("API correctly rejected empty email: %s", resp.Message)
	}
}

func TestIntegration_GetKycInformation_SpecialCharactersInEmail(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)

	testEmails := []string{
		"test+tag@example.com",
		"test.name@example.com",
		"test_name@example.com",
	}

	for _, email := range testEmails {
		t.Run(email, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.KycInformationRequest{
				UserEmail: email,
			}

			resp, err := client.GetKycInformation(ctx, req)
			if err != nil {
				t.Logf("GetKycInformation returned error for email %s: %v", email, err)

				return
			}

			t.Logf("Response for %s: success=%v", email, resp.Success)
		})
	}
}

func TestIntegration_GetKycInformation_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	}

	_, err := client.GetKycInformation(ctx, req)
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

func TestIntegration_GetKycInformation_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.KycInformationRequest{
		UserEmail: "test@example.com",
	}

	cancel()

	_, err := client.GetKycInformation(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}

func TestIntegration_GetKycInformation_MultipleRequests(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	emails := getTestEmails(t)

	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.KycInformationRequest{
				UserEmail: email,
			}

			resp, err := client.GetKycInformation(ctx, req)
			if err != nil {
				t.Logf("GetKycInformation error for %s: %v", email, err)

				return
			}

			t.Logf("Response for %s: success=%v, status=%s",
				email, resp.Success, resp.Data.KycStatus)
		})
	}
}

func TestIntegration_GetKycInformation_LongTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: getTestEmail(t),
	}

	start := time.Now()
	resp, err := client.GetKycInformation(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	t.Logf("Request completed in %v", duration)
	t.Logf("Response: success=%v, status=%s", resp.Success, resp.Data.KycStatus)
}

func TestIntegration_GetKycInformation_ResponseDataFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: getTestEmail(t),
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful, skipping field validation: %s", resp.Message)
	}

	data := resp.Data

	// Log all fields for debugging
	t.Logf("Response Data Fields:")
	t.Logf("  FirstName: %q", data.FirstName)
	t.Logf("  LastName: %q", data.LastName)
	t.Logf("  DateOfBirth: %q", data.DateOfBirth)
	t.Logf("  Gender: %q", data.Gender)
	t.Logf("  Nationality: %q", data.Nationality)
	t.Logf("  IDType: %q", data.IDType)
	t.Logf("  NationalID: %q", data.NationalID)
	t.Logf("  IssueDate: %q", data.IssueDate)
	t.Logf("  ExpiryDate: %q", data.ExpiryDate)
	t.Logf("  Address: %q", data.Address)
	t.Logf("  PhoneNumber: %q", data.PhoneNumber)
	t.Logf("  PhoneCountryCode: %q", data.PhoneCountryCode)
	t.Logf("  KycStatus: %q", data.KycStatus)
	t.Logf("  RejectReason: %q", data.RejectReason)
	t.Logf("  FrontIDImage length: %d", len(data.FrontIDImage))
	t.Logf("  BackIDImage length: %d", len(data.BackIDImage))
	t.Logf("  HoldIDImage length: %d", len(data.HoldIDImage))

	if data.KycStatus == "" {
		t.Error("KycStatus should not be empty")
	}
}

func TestIntegration_GetKycInformation_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	emails := getTestEmails(t)

	results := make(chan struct {
		email   string
		success bool
		err     error
	}, len(emails))

	for _, email := range emails {
		go func() {
			req := &goaliniex.KycInformationRequest{
				UserEmail: email,
			}

			resp, err := client.GetKycInformation(ctx, req)

			result := struct {
				email   string
				success bool
				err     error
			}{
				email:   email,
				success: false,
				err:     err,
			}

			if err == nil && resp != nil {
				result.success = resp.Success
			}

			results <- result
		}()
	}

	for range emails {
		result := <-results
		if result.err != nil {
			t.Logf("Concurrent request for %s failed: %v", result.email, result.err)
		} else {
			t.Logf("Concurrent request for %s: success=%v", result.email, result.success)
		}
	}
}

func TestIntegration_GetKycInformation_RepeatedRequests(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: getTestEmail(t),
	}

	for requestNum := range 3 {
		resp, err := client.GetKycInformation(ctx, req)
		if err != nil {
			t.Errorf("Request %d failed: %v", requestNum+1, err)

			continue
		}

		if !resp.Success {
			t.Errorf("Request %d returned success=false", requestNum+1)

			continue
		}

		t.Logf("Request %d: success=%v, status=%s", requestNum+1, resp.Success, resp.Data.KycStatus)
	}
}
