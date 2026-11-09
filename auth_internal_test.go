package goaliniex

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateAuthTestKeyPair(t *testing.T) ([]byte, *rsa.PublicKey) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return privateKeyPEM, &privateKey.PublicKey
}

func TestBuildClientAssertion_ClaimsAndSignature(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKey := generateAuthTestKeyPair(t)

	assertion, err := buildClientAssertion(privateKeyPEM, "test-partner", time.Minute)
	if err != nil {
		t.Fatalf("buildClientAssertion: %v", err)
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(assertion, &claims, func(_ *jwt.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		t.Fatalf("parse/verify assertion: %v", err)
	}

	if !token.Valid {
		t.Fatal("expected token to be valid")
	}

	if claims["iss"] != "test-partner" {
		t.Errorf("iss = %v, want test-partner", claims["iss"])
	}

	if claims["aud"] != clientAssertionAudience {
		t.Errorf("aud = %v, want %v", claims["aud"], clientAssertionAudience)
	}

	for _, key := range []string{"iat", "exp", "jti"} {
		if _, ok := claims[key]; !ok {
			t.Errorf("missing claim %q", key)
		}
	}
}

func TestBuildClientAssertion_UniqueJTI(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateAuthTestKeyPair(t)

	a1, err := buildClientAssertion(privateKeyPEM, "test-partner", time.Minute)
	if err != nil {
		t.Fatalf("buildClientAssertion: %v", err)
	}

	a2, err := buildClientAssertion(privateKeyPEM, "test-partner", time.Minute)
	if err != nil {
		t.Fatalf("buildClientAssertion: %v", err)
	}

	if a1 == a2 {
		t.Error("expected two distinct assertions (different jti/iat), got identical tokens")
	}
}

// countingHTTPClient returns one canned response per call, in order, and
// records how many calls were made.
type countingHTTPClient struct {
	responses []*http.Response
	calls     int
}

func (c *countingHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	resp := c.responses[c.calls]
	c.calls++

	return resp, nil
}

func newCannedBodyResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func tokenResponseBody(accessToken string, expiresAt time.Time) string {
	return `{"code":"SUCCESS","message":"ok","data":{"accessToken":"` + accessToken +
		`","tokenType":"Bearer","expiresAt":` + strconv.FormatInt(expiresAt.Unix(), 10) + `}}`
}

func TestTokenCache_HitAvoidsSecondAuthCall(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateAuthTestKeyPair(t)

	mock := &countingHTTPClient{
		responses: []*http.Response{
			newCannedBodyResponse(http.StatusOK, tokenResponseBody("tok-1", time.Now().Add(time.Hour))),
		},
	}

	client, err := NewClient("https://sandbox.alixpay.com", "test-partner", privateKeyPEM, WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()

	tok1, err := client.tokenCache.getToken(ctx, client, "user@example.com")
	if err != nil {
		t.Fatalf("getToken (1st): %v", err)
	}

	tok2, err := client.tokenCache.getToken(ctx, client, "user@example.com")
	if err != nil {
		t.Fatalf("getToken (2nd): %v", err)
	}

	if tok1 != tok2 {
		t.Errorf("expected cached token to be reused, got %q then %q", tok1, tok2)
	}

	if mock.calls != 1 {
		t.Errorf("expected exactly 1 auth call (cache hit on 2nd), got %d", mock.calls)
	}
}

func TestTokenCache_RenewsNearExpiry(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateAuthTestKeyPair(t)

	// expiresAt is already within the renew skew below, forcing a renewal
	// on every call.
	soonBody := tokenResponseBody("tok-soon", time.Now().Add(time.Second))

	mock := &countingHTTPClient{
		responses: []*http.Response{
			newCannedBodyResponse(http.StatusOK, soonBody),
			newCannedBodyResponse(http.StatusOK, soonBody),
		},
	}

	client, err := NewClient(
		"https://sandbox.alixpay.com", "test-partner", privateKeyPEM,
		WithHTTPClient(mock), WithTokenRenewSkew(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()

	if _, err := client.tokenCache.getToken(ctx, client, "user@example.com"); err != nil {
		t.Fatalf("getToken (1st): %v", err)
	}

	if _, err := client.tokenCache.getToken(ctx, client, "user@example.com"); err != nil {
		t.Fatalf("getToken (2nd): %v", err)
	}

	if mock.calls != 2 {
		t.Errorf("expected 2 auth calls (renewal forced each time), got %d", mock.calls)
	}
}

func TestTokenCache_IsolatedPerUserEmail(t *testing.T) {
	t.Parallel()

	privateKeyPEM, _ := generateAuthTestKeyPair(t)

	mock := &countingHTTPClient{
		responses: []*http.Response{
			newCannedBodyResponse(http.StatusOK, tokenResponseBody("tok-a", time.Now().Add(time.Hour))),
			newCannedBodyResponse(http.StatusOK, tokenResponseBody("tok-b", time.Now().Add(time.Hour))),
		},
	}

	client, err := NewClient("https://sandbox.alixpay.com", "test-partner", privateKeyPEM, WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()

	tokA, err := client.tokenCache.getToken(ctx, client, "a@example.com")
	if err != nil {
		t.Fatalf("getToken a: %v", err)
	}

	tokB, err := client.tokenCache.getToken(ctx, client, "b@example.com")
	if err != nil {
		t.Fatalf("getToken b: %v", err)
	}

	if tokA == tokB {
		t.Errorf("expected different tokens per user email, got %q for both", tokA)
	}

	if mock.calls != 2 {
		t.Errorf("expected 2 auth calls (one per user), got %d", mock.calls)
	}
}
