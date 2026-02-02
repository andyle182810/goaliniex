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
	ErrInvalidPEM        = errors.New("invalid private key PEM")
	ErrInvalidPublicPEM  = errors.New("invalid public key PEM")
	ErrUnsupportedKey    = errors.New("unsupported private key type")
	ErrUnsupportedPubKey = errors.New("unsupported public key type")
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

func Verify(publicKeyPEM []byte, payload []byte, signatureBase64 string) error {
	rsaPublicKey, err := parseRSAPublicKey(publicKeyPEM)
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

func parseRSAPublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
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
