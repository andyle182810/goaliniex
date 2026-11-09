package signer

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrInvalidPEM        = errors.New("invalid private key PEM")
	ErrInvalidPublicPEM  = errors.New("invalid public key PEM")
	ErrUnsupportedKey    = errors.New("unsupported private key type")
	ErrUnsupportedPubKey = errors.New("unsupported public key type")
)

// ParsePrivateKey parses an RSA private key from a PEM block, supporting both
// PKCS#1 and PKCS#8 encodings.
func ParsePrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, ErrInvalidPEM
	}

	// PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrUnsupportedKey
	}

	return rsaKey, nil
}

// Verify checks an RSA-SHA256 (PKCS#1 v1.5) signature, base64-standard-encoded,
// over payload using the given PEM-encoded public key. Used to verify Alix
// webhook signatures (X-Alix-Signature).
func Verify(publicKeyPEM []byte, payload []byte, signatureBase64 string) error {
	rsaPublicKey, err := ParsePublicKey(publicKeyPEM)
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	digest := sha256.Sum256(payload)

	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, digest[:], signature)
}

// ParsePublicKey parses an RSA public key from a PEM block, supporting PKIX,
// PKCS#1, and X.509 certificate encodings.
func ParsePublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, ErrInvalidPublicPEM
	}

	// Try PKIX format (most common for public keys)
	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaKey, ok := pub.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}

		return nil, ErrUnsupportedPubKey
	}

	// Try PKCS#1 format
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try certificate format
	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		if rsaKey, ok := cert.PublicKey.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}

		return nil, ErrUnsupportedPubKey
	}

	return nil, fmt.Errorf("parse public key: %w", ErrInvalidPublicPEM)
}
