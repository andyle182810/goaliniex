package goaliniex_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andyle182810/goaliniex/v2"
)

func generateWebhookTestKeyPair(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})

	return privateKeyPEM, publicKeyPEM
}

func signWebhookPayload(t *testing.T, privateKeyPEM []byte, timestamp, body string) string {
	t.Helper()

	block, _ := pem.Decode(privateKeyPEM)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	digest := sha256.Sum256([]byte(timestamp + "." + body))

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature)
}

func TestVerifyWebhookSignature_Valid(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateWebhookTestKeyPair(t)
	timestamp := "1700000000"
	body := `{"event":"order.success","orderId":"abc"}`

	signature := signWebhookPayload(t, privateKeyPEM, timestamp, body)

	err := goaliniex.VerifyWebhookSignature(publicKeyPEM, timestamp, []byte(body), signature)
	if err != nil {
		t.Errorf("expected valid signature to verify, got: %v", err)
	}
}

func TestVerifyWebhookSignature_TamperedBody(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateWebhookTestKeyPair(t)
	timestamp := "1700000000"
	body := `{"event":"order.success","orderId":"abc"}`

	signature := signWebhookPayload(t, privateKeyPEM, timestamp, body)

	err := goaliniex.VerifyWebhookSignature(publicKeyPEM, timestamp, []byte(`{"event":"order.failed","orderId":"abc"}`), signature)
	if !errors.Is(err, goaliniex.ErrInvalidWebhookSignature) {
		t.Errorf("expected ErrInvalidWebhookSignature for tampered body, got: %v", err)
	}
}

func TestVerifyWebhookSignature_TamperedTimestamp(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateWebhookTestKeyPair(t)
	body := `{"event":"order.success","orderId":"abc"}`

	signature := signWebhookPayload(t, privateKeyPEM, "1700000000", body)

	err := goaliniex.VerifyWebhookSignature(publicKeyPEM, "1700000001", []byte(body), signature)
	if !errors.Is(err, goaliniex.ErrInvalidWebhookSignature) {
		t.Errorf("expected ErrInvalidWebhookSignature for tampered timestamp, got: %v", err)
	}
}

func TestVerifyWebhookSignature_WrongKey(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateWebhookTestKeyPair(t)
	_, wrongPublicKeyPEM := generateWebhookTestKeyPair(t)
	timestamp := "1700000000"
	body := `{"event":"order.success","orderId":"abc"}`

	signature := signWebhookPayload(t, privateKeyPEM, timestamp, body)

	err := goaliniex.VerifyWebhookSignature(wrongPublicKeyPEM, timestamp, []byte(body), signature)
	if !errors.Is(err, goaliniex.ErrInvalidWebhookSignature) {
		t.Errorf("expected ErrInvalidWebhookSignature for wrong key, got: %v", err)
	}
}

func TestVerifyWebhookSignature_GarbageBase64(t *testing.T) {
	t.Parallel()

	_, publicKeyPEM := generateWebhookTestKeyPair(t)

	err := goaliniex.VerifyWebhookSignature(publicKeyPEM, "1700000000", []byte("body"), "not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for garbage base64 signature")
	}
}

func TestVerifyWebhookRequest_RestoresBody(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateWebhookTestKeyPair(t)
	timestamp := "1700000000"
	body := `{"event":"order.success","orderId":"abc"}`

	signature := signWebhookPayload(t, privateKeyPEM, timestamp, body)

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set("X-Alix-Signature", signature)
	req.Header.Set("X-Alix-Timestamp", timestamp)

	gotBody, err := goaliniex.VerifyWebhookRequest(publicKeyPEM, req)
	if err != nil {
		t.Fatalf("VerifyWebhookRequest: %v", err)
	}

	if string(gotBody) != body {
		t.Errorf("returned body = %q, want %q", gotBody, body)
	}

	// Body must still be readable by a downstream handler.
	remaining, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("read restored body: %v", err)
	}

	if string(remaining) != body {
		t.Errorf("restored body = %q, want %q", remaining, body)
	}
}
