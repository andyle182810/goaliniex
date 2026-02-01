package goaliniex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/andyle182810/goaliniex/signer"
)

const UserAgent = "aliniex-go-sdk"

var (
	// Client initialization errors.
	ErrEmptyBaseURL     = errors.New("baseURL is required")
	ErrEmptyPartnerCode = errors.New("partnerCode is required")
	ErrEmptySecretKey   = errors.New("secretKey is required")
	ErrEmptyPrivateKey  = errors.New("privateKey is required")

	// Request lifecycle errors.
	ErrNilRequest    = errors.New("request is nil")
	ErrRequestBuild  = errors.New("failed to build request")
	ErrRequestSign   = errors.New("failed to sign request")
	ErrRequestEncode = errors.New("failed to encode request body")
	ErrInvalidParams = errors.New("invalid request params")

	// HTTP / transport errors.
	ErrHTTPFailure      = errors.New("http request failed")
	ErrUnexpectedStatus = errors.New("unexpected http status code")
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(*Client)

type Client struct {
	baseURL     string
	partnerCode string
	secretKey   string
	privateKey  []byte
	logger      Logger
	debug       bool
	httpClient  HTTPClient
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func NewClient(
	baseURL string,
	partnerCode string,
	secretKey string,
	privateKey []byte,
	opts ...Option,
) (*Client, error) {
	if baseURL == "" {
		return nil, ErrEmptyBaseURL
	}

	if partnerCode == "" {
		return nil, ErrEmptyPartnerCode
	}

	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}

	if len(privateKey) == 0 {
		return nil, ErrEmptyPrivateKey
	}

	client := &Client{
		baseURL:     baseURL,
		partnerCode: partnerCode,
		secretKey:   secretKey,
		privateKey:  privateKey,
		httpClient:  http.DefaultClient,
		logger:      slog.Default(),
		debug:       false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) logDebug(msg string, attrs ...any) {
	if c.debug {
		c.logger.Debug(msg, attrs...)
	}
}

func paramsToMap(params any) (map[string]any, error) {
	if params == nil {
		return map[string]any{}, nil
	}

	if m, ok := params.(map[string]any); ok {
		return m, nil
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) buildGETRequest(req *request, fullURL string, headers http.Header) error {
	paramsMap, err := paramsToMap(req.Params)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidParams, err)
	}

	if len(paramsMap) > 0 {
		queryParams := url.Values{}
		for key, value := range paramsMap {
			queryParams.Add(key, fmt.Sprintf("%v", value))
		}

		fullURL += "?" + queryParams.Encode()
	}

	c.logDebug("http request", "url", fullURL)

	req.FullURL = fullURL
	req.Header = headers
	req.Body = nil

	return nil
}

func (c *Client) buildRequest(req *request) error {
	if req == nil {
		return ErrNilRequest
	}

	fullURL := c.baseURL + req.Endpoint

	headers := http.Header{}
	if req.Header != nil {
		headers = req.Header.Clone()
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", UserAgent)

	// For GET requests, use query parameters instead of body
	if req.Method == http.MethodGet {
		return c.buildGETRequest(req, fullURL, headers)
	}

	// For POST/other methods, use JSON body
	bodyMap, err := paramsToMap(req.Params)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidParams, err)
	}

	if !req.Public {
		signature, err := signer.Sign(c.privateKey, req.SigningData)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRequestSign, err)
		}

		bodyMap["partnerCode"] = c.partnerCode
		bodyMap["signature"] = signature
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestEncode, err)
	}

	c.logDebug("http request", "url", fullURL)
	c.logDebug("http request body", "body", string(bodyBytes))

	req.FullURL = fullURL
	req.Header = headers
	req.Body = bytes.NewReader(bodyBytes)

	return nil
}

func (c *Client) execute(ctx context.Context, req *request) ([]byte, error) {
	if err := c.buildRequest(req); err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		req.FullURL,
		req.Body,
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header = req.Header

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTPFailure, err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.logDebug("http response", "status", resp.StatusCode)
	c.logDebug("http response body", "body", string(responseBody))

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"%w: status=%d body=%s",
			ErrUnexpectedStatus,
			resp.StatusCode,
			string(responseBody),
		)
	}

	return responseBody, nil
}
