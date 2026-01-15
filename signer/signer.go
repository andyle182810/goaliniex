package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrInvalidPEM     = errors.New("invalid private key PEM")
	ErrUnsupportedKey = errors.New("unsupported private key type")
)

func Sign(privateKeyPEM []byte, payload []byte) (string, error) {
	rsaPrivateKey, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	digest := sha256.Sum256(payload)

	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		rsaPrivateKey,
		crypto.SHA256,
		digest[:],
	)
	if err != nil {
		return "", fmt.Errorf("sign payload: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func parseRSAPrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
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
