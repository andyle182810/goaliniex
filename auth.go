package goaliniex

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/andyle182810/goaliniex/v2/signer"
)

// clientAssertionAudience is the fixed "aud" claim Alix expects on the
// RS256 client-assertion JWT exchanged for an access token.
const clientAssertionAudience = "alixpay-v3-token"

// TokenResponse is the data payload of POST /v3/auth/tokens.
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type tokenRequest struct {
	ClientAssertion string `json:"clientAssertion"`
	UserEmail       string `json:"userEmail"`
}

// buildClientAssertion signs a short-lived RS256 JWT identifying the partner,
// per Alix v3's auth handshake (POST /v3/auth/tokens).
func buildClientAssertion(privateKeyPEM []byte, partnerCode string, ttl time.Duration) (string, error) {
	rsaPrivateKey, err := signer.ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	jti, err := randomJTI()
	if err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"iss": partnerCode,
		"aud": clientAssertionAudience,
		"iat": now.Unix(),
		"jti": jti,
		"exp": now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign client assertion: %w", err)
	}

	return signed, nil
}

func randomJTI() (string, error) {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

// requestToken exchanges a freshly-signed client assertion for a bearer
// access token scoped to userEmail, via the open POST /v3/auth/tokens
// endpoint.
func (c *Client) requestToken(ctx context.Context, userEmail string) (*TokenResponse, error) {
	assertion, err := buildClientAssertion(c.privateKey, c.partnerCode, c.jwtTTL)
	if err != nil {
		return nil, fmt.Errorf("build client assertion: %w", err)
	}

	apiRequest := request{
		Method:   http.MethodPost,
		Endpoint: "/v3/auth/tokens",
		Params: &tokenRequest{
			ClientAssertion: assertion,
			UserEmail:       userEmail,
		},
		Public: true,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[TokenResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	if !response.IsSuccess() || response.Data == nil {
		return nil, fmt.Errorf("%w: code=%s message=%s", ErrAuthFailed, response.Code, response.Message)
	}

	return response.Data, nil
}

type cachedToken struct {
	accessToken string
	expiresAt   time.Time
}

// tokenCache holds bearer access tokens keyed by user email (or the empty
// string / partner code for partner-level, non-user-scoped endpoints), and
// transparently renews them shortly before they expire.
type tokenCache struct {
	mu     sync.Mutex
	tokens map[string]*cachedToken
}

func newTokenCache() *tokenCache {
	return &tokenCache{
		tokens: make(map[string]*cachedToken),
	}
}

func (tc *tokenCache) getToken(ctx context.Context, client *Client, userEmail string) (string, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if cached, ok := tc.tokens[userEmail]; ok {
		if time.Now().Add(client.tokenRenewSkew).Before(cached.expiresAt) {
			return cached.accessToken, nil
		}
	}

	tokenResp, err := client.requestToken(ctx, userEmail)
	if err != nil {
		return "", err
	}

	cached := &cachedToken{
		accessToken: tokenResp.AccessToken,
		expiresAt:   time.Unix(tokenResp.ExpiresAt, 0),
	}

	tc.tokens[userEmail] = cached

	return cached.accessToken, nil
}
