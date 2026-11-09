package goaliniex_test

import (
	"context"
	"testing"
	"time"
)

func TestIntegration_GetPartnerFund(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.GetPartnerFund(ctx, "")
	if err != nil {
		t.Fatalf("GetPartnerFund returned error: %v", err)
	}

	t.Logf("GetPartnerFund: code=%s message=%s", resp.Code, resp.Message)

	if resp.IsSuccess() && resp.Data != nil {
		t.Logf("cryptoCurrency=%s balance=%.4f depositAddresses=%d",
			resp.Data.CryptoCurrency, resp.Data.Balance, len(resp.Data.DepositAddresses))
	}
}
