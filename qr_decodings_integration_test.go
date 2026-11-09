package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex/v2"
)

// DecodeQR is an open (public) endpoint, so it needs neither a partner code
// nor a private key — only the base sandbox client.
func TestIntegration_DecodeQR_UnsupportedInput(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.DecodeQR(ctx, &goaliniex.DecodeQRRequest{QR: "not-a-real-qr-payload"})
	if err != nil {
		t.Logf("DecodeQR returned error (may be expected for invalid input): %v", err)

		return
	}

	t.Logf("DecodeQR: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("supported=%v qrType=%s countryCode=%s", resp.Data.Supported, resp.Data.QRType, resp.Data.CountryCode)
	}
}
