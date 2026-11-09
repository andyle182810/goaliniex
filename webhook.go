package goaliniex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andyle182810/goaliniex/v2/signer"
)

// ErrInvalidWebhookSignature is returned when an incoming webhook's
// X-Alix-Signature does not verify against the given Alix public key.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

// WebhookEvent is the order-lifecycle callback body Alix POSTs when an
// order reaches a terminal state (order.success / order.failed). Delivery
// is once, with no retries: handlers must be idempotent by OrderID and
// treat GET /v3/{service}/orders/{orderId} as the source of truth if a
// delivery is missed.
type WebhookEvent struct {
	Event      string `json:"event"`
	Service    string `json:"service"`
	OrderID    string `json:"orderId"`
	OccurredAt string `json:"occurredAt"`
	// Data is the full order snapshot at event time, same shape as
	// GET /v3/{service}/orders/{orderId}; unmarshal it into the order type
	// matching Service (e.g. ScanToPayOrder, CashoutOrder).
	Data json.RawMessage `json:"data"`
}

// VerifyWebhookSignature verifies an Alix webhook's X-Alix-Signature:
// base64-standard RSA-SHA256 over "<timestamp>.<rawBody>", using Alix's
// public key (issued to the partner out-of-band, distinct from the
// partner's own key pair used for JWT auth).
func VerifyWebhookSignature(alixPublicKeyPEM []byte, timestamp string, rawBody []byte, signatureBase64 string) error {
	payload := append([]byte(timestamp+"."), rawBody...)

	if err := signer.Verify(alixPublicKeyPEM, payload, signatureBase64); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidWebhookSignature, err)
	}

	return nil
}

// VerifyWebhookRequest reads and verifies the X-Alix-Signature/
// X-Alix-Timestamp headers of an incoming webhook *http.Request, restores
// r.Body so downstream handlers can still read it, and returns the raw body
// for the caller to unmarshal into a WebhookEvent.
func VerifyWebhookRequest(alixPublicKeyPEM []byte, r *http.Request) ([]byte, error) {
	signature := r.Header.Get("X-Alix-Signature")
	timestamp := r.Header.Get("X-Alix-Timestamp")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read webhook body: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	if err := VerifyWebhookSignature(alixPublicKeyPEM, timestamp, body, signature); err != nil {
		return nil, err
	}

	return body, nil
}
