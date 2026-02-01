package goaliniex_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

func TestIntegration_GetQRCodeInfo_ValidVietQR(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
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

	t.Logf("QR Code Information retrieved successfully")
	t.Logf("  Bank Account Number: %s", resp.Data.BankAccountNumber)
	t.Logf("  Bank Code: %s", resp.Data.BankCode)
	t.Logf("  Bank Name: %s", resp.Data.BankName)
	t.Logf("  Country Code: %s", resp.Data.CountryCode)
	t.Logf("  QR Type: %s", resp.Data.QRType)
	t.Logf("  Amount: %.2f", resp.Data.Amount)
	t.Logf("  Additional Data: %+v", resp.Data.AdditionalData)
}

func TestIntegration_GetQRCodeInfo_InvalidQRContent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "invalid-qr-content-12345",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Logf("GetQRCodeInfo returned error (may be expected): %v", err)

		return
	}

	if !resp.Success {
		t.Logf("API correctly indicated invalid QR code: %s (error code: %d)", resp.Message, resp.ErrorCode)
	} else {
		t.Logf("Warning: API accepted invalid QR content, this might be unexpected")
	}
}

func TestIntegration_GetQRCodeInfo_EmptyQRContent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Logf("GetQRCodeInfo returned error for empty QR content: %v", err)

		return
	}

	if resp.Success {
		t.Logf("Warning: API accepted empty QR content, this might be a validation issue")
	} else {
		t.Logf("API correctly rejected empty QR content: %s", resp.Message)
	}
}

func TestIntegration_GetQRCodeInfo_ContextTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	_, err := client.GetQRCodeInfo(ctx, req)
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

func TestIntegration_GetQRCodeInfo_ContextCancellation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	cancel()

	_, err := client.GetQRCodeInfo(ctx, req)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	} else {
		t.Logf("Context cancellation error as expected: %v", err)
	}
}

func TestIntegration_GetQRCodeInfo_MultipleQRFormats(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	// Different QR code formats to test
	qrCodes := map[string]string{
		"VietQR": "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
		// Add more QR formats when available
	}

	for name, qrContent := range qrCodes {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.GetQRCodeInfoRequest{
				QRContent: qrContent,
			}

			resp, err := client.GetQRCodeInfo(ctx, req)
			if err != nil {
				t.Logf("GetQRCodeInfo error for %s: %v", name, err)

				return
			}

			t.Logf("Response for %s: success=%v, type=%s, country=%s",
				name, resp.Success, resp.Data.QRType, resp.Data.CountryCode)
		})
	}
}

func TestIntegration_GetQRCodeInfo_ResponseDataFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if !resp.Success {
		t.Skipf("Request was not successful, skipping field validation: %s", resp.Message)
	}

	data := resp.Data

	// Log all fields for debugging
	t.Logf("Response Data Fields:")
	t.Logf("  BankAccountNumber: %q", data.BankAccountNumber)
	t.Logf("  BankCode: %q", data.BankCode)
	t.Logf("  BankName: %q", data.BankName)
	t.Logf("  CountryCode: %q", data.CountryCode)
	t.Logf("  QRType: %q", data.QRType)
	t.Logf("  Amount: %.2f", data.Amount)
	t.Logf("  AdditionalData keys: %d", len(data.AdditionalData))

	for key, value := range data.AdditionalData {
		t.Logf("    %s: %v", key, value)
	}

	if data.BankAccountNumber == "" {
		t.Error("BankAccountNumber should not be empty")
	}

	if data.CountryCode == "" {
		t.Error("CountryCode should not be empty")
	}

	if data.QRType == "" {
		t.Error("QRType should not be empty")
	}
}

func TestIntegration_GetQRCodeInfo_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	qrCodes := []string{
		"00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
		"00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	results := make(chan struct {
		qrContent string
		success   bool
		err       error
	}, len(qrCodes))

	for _, qrContent := range qrCodes {
		go func() {
			req := &goaliniex.GetQRCodeInfoRequest{
				QRContent: qrContent,
			}

			resp, err := client.GetQRCodeInfo(ctx, req)

			result := struct {
				qrContent string
				success   bool
				err       error
			}{
				qrContent: qrContent[:20] + "...",
				success:   false,
				err:       err,
			}

			if err == nil && resp != nil {
				result.success = resp.Success
			}

			results <- result
		}()
	}

	for range qrCodes {
		result := <-results
		if result.err != nil {
			t.Logf("Concurrent request for QR %s failed: %v", result.qrContent, result.err)
		} else {
			t.Logf("Concurrent request for QR %s: success=%v", result.qrContent, result.success)
		}
	}
}

func TestIntegration_GetQRCodeInfo_RepeatedRequests(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	for requestNum := range 3 {
		resp, err := client.GetQRCodeInfo(ctx, req)
		if err != nil {
			t.Errorf("Request %d failed: %v", requestNum+1, err)

			continue
		}

		if !resp.Success {
			t.Errorf("Request %d returned success=false: %s", requestNum+1, resp.Message)

			continue
		}

		t.Logf("Request %d: success=%v, qrType=%s, country=%s",
			requestNum+1, resp.Success, resp.Data.QRType, resp.Data.CountryCode)
	}
}

func TestIntegration_GetQRCodeInfo_LongTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	start := time.Now()
	resp, err := client.GetQRCodeInfo(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	t.Logf("Request completed in %v", duration)
	t.Logf("Response: success=%v, qrType=%s, country=%s",
		resp.Success, resp.Data.QRType, resp.Data.CountryCode)
}

func TestIntegration_GetQRCodeInfo_MalformedQRContent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)

	testCases := []struct {
		name      string
		qrContent string
	}{
		{
			name:      "Special characters",
			qrContent: "!@#$%^&*()",
		},
		{
			name:      "Very long string",
			qrContent: strings.Repeat("A", 1000),
		},
		{
			name:      "Only numbers",
			qrContent: "1234567890",
		},
		{
			name:      "Whitespace",
			qrContent: "   ",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			req := &goaliniex.GetQRCodeInfoRequest{
				QRContent: testCase.qrContent,
			}

			resp, err := client.GetQRCodeInfo(ctx, req)
			if err != nil {
				t.Logf("GetQRCodeInfo returned error for %s: %v", testCase.name, err)

				return
			}

			t.Logf("Response for %s: success=%v, message=%s",
				testCase.name, resp.Success, resp.Message)
		})
	}
}
