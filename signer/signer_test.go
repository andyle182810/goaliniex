package signer_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/andyle182810/goaliniex/signer"
)

func generateTestKeyPair(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM
}

func generatePKCS1PublicKey(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return privateKeyPEM, publicKeyPEM
}

func TestVerify_ValidSignature(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	payload := []byte("test payload data")

	signature, err := signer.Sign(privateKeyPEM, payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	err = signer.Verify(publicKeyPEM, payload, signature)
	if err != nil {
		t.Errorf("verify valid signature: %v", err)
	}
}

func TestVerify_InvalidSignature(t *testing.T) {
	t.Parallel()

	_, publicKeyPEM := generateTestKeyPair(t)
	payload := []byte("test payload data")

	err := signer.Verify(publicKeyPEM, payload, "aW52YWxpZHNpZ25hdHVyZQ==")
	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestVerify_TamperedPayload(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	payload := []byte("original payload")

	signature, err := signer.Sign(privateKeyPEM, payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	tamperedPayload := []byte("tampered payload")

	err = signer.Verify(publicKeyPEM, tamperedPayload, signature)
	if err == nil {
		t.Error("expected error for tampered payload")
	}
}

func TestVerify_InvalidBase64Signature(t *testing.T) {
	t.Parallel()

	_, publicKeyPEM := generateTestKeyPair(t)
	payload := []byte("test payload")

	err := signer.Verify(publicKeyPEM, payload, "not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestVerify_InvalidPublicKeyPEM(t *testing.T) {
	t.Parallel()

	payload := []byte("test payload")

	err := signer.Verify([]byte("invalid pem"), payload, "c2lnbmF0dXJl")
	if !errors.Is(err, signer.ErrInvalidPublicPEM) {
		t.Errorf("expected ErrInvalidPublicPEM, got: %v", err)
	}
}

func TestVerify_WrongPublicKey(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateTestKeyPair(t)
	_, wrongPublicKeyPEM := generateTestKeyPair(t)
	payload := []byte("test payload")

	signature, err := signer.Sign(privateKeyPEM, payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	err = signer.Verify(wrongPublicKeyPEM, payload, signature)
	if err == nil {
		t.Error("expected error when verifying with wrong public key")
	}
}

func TestVerify_PKCS1PublicKey(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generatePKCS1PublicKey(t)
	payload := []byte("test payload for PKCS1")

	signature, err := signer.Sign(privateKeyPEM, payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	err = signer.Verify(publicKeyPEM, payload, signature)
	if err != nil {
		t.Errorf("verify with PKCS1 public key: %v", err)
	}
}

func TestVerify_EmptyPayload(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	payload := []byte{}

	signature, err := signer.Sign(privateKeyPEM, payload)
	if err != nil {
		t.Fatalf("sign empty payload: %v", err)
	}

	err = signer.Verify(publicKeyPEM, payload, signature)
	if err != nil {
		t.Errorf("verify empty payload: %v", err)
	}
}

func TestVerify_FromFile(t *testing.T) {
	t.Parallel()

	publicKeyPath := os.Getenv("TEST_PUBLIC_KEY_PATH")
	signature := os.Getenv("TEST_SIGNATURE")
	payload := os.Getenv("TEST_PAYLOAD")

	if publicKeyPath == "" || signature == "" || payload == "" {
		t.Skip("Set TEST_PUBLIC_KEY_PATH, TEST_SIGNATURE, and TEST_PAYLOAD to run this test")
	}

	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Fatalf("read public key file: %v", err)
	}

	t.Logf("Payload (%d bytes): %q", len(payload), payload)
	t.Logf("Payload hex: %x", []byte(payload))
	t.Logf("Signature length: %d", len(signature))

	err = signer.Verify(publicKeyPEM, []byte(payload), signature)
	if err != nil {
		t.Errorf("verify signature: %v", err)
	} else {
		t.Log("Signature verified successfully")
	}
}

func TestSignAndVerify_RoundTrip(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)

	testCases := []struct {
		name    string
		payload []byte
	}{
		{"simple text", []byte("hello world")},
		{"json data", []byte(`{"key": "value", "number": 123}`)},
		{"binary data", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}},
		{"large payload", make([]byte, 10000)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			signature, err := signer.Sign(privateKeyPEM, testCase.payload)
			if err != nil {
				t.Fatalf("sign: %v", err)
			}

			err = signer.Verify(publicKeyPEM, testCase.payload, signature)
			if err != nil {
				t.Errorf("verify: %v", err)
			}
		})
	}
}
