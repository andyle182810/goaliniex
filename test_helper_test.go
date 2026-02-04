package goaliniex_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		Status:           http.StatusText(statusCode),
		StatusCode:       statusCode,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           make(http.Header),
		Body:             io.NopCloser(bytes.NewBufferString(body)),
		ContentLength:    int64(len(body)),
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
}

func testPrivateKey() []byte {
	return []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDQoOm2Yf7m63+Q
PTVn9Vf9lXeSNZ3DJQjlQ7WMmmZ/ytgp/UFUnkwtZ5IP3g73Xg8xMXftYqQl+7wa
4a15Y5SF6Ju6FIveBNxLd1Q5U/dFf2ctSKFB/OLvOBk0o+F85PHncf0AOT9dFjhp
gErWZ6i686Z3gCL3sc6CZx5pAQq54ovs98UIUMXbrLUgn618iN8S+bqohEQG7FJi
DsLY1EMf/FmwP/kuoX+98ymCbXR+7KBCz9dHmfK+bb2qHwL7ky/VwFwNI/i+nDBi
upqLMpKsgiFOUiH0NfgmAy61ywc7wZXPgKax2edTreUItxMP9p5eyRuTZ0rAX/Xv
Vvj3OOrPAgMBAAECggEAZo9sRTIIhsmnlKdabUgxoOJM/S3pZ+j5ZgdypPO0RsdX
m9C5wJjvFvUO20kNL3LoYgURETxIOLn3f3mI959KAMhYYdI+7f6Ss3CukV4jNWGL
vbmyOIfSocoArh2QOH8uIlqphuYdravCQK8GWz9rNbiYka+GwSCCRh3eV71g4f7G
2M0F2dgfNCR8LvnLImxQN9MNkyIQTRygRh3h5vlSjDoF+l7TEUHt9E6iyrJ7t7MO
3vQvOqdWFU1eHIEK15OKxjEke1MCOyJjVPFsJummuK6w8xPMH7OQRvkAYhyn7rUx
nxPhaYNa1fC+ImyTEW//aJySc00TIdV2U1Wbj4AsWQKBgQDx4KaAj/TU7WrQSUR+
EBnwNiZNLxRuEwekMJSfpqQirm2qJFuDd4zXfAWju6aeAUvoz7d6rbIZhSxh9umj
8WxpTjUzgBNiEpWzdDn/rMdqRi9xgPKBkmSI4Skyxhd50+v3mHVaMfwqOTrwwRQa
THLjo6ij0bYh8oc/s/OxjiuvmQKBgQDcz0oPKl+YjaL4stYrIIWO4sZ49XTxCw1t
q32T9G0Nba3dOpMgCW7Xs3/JofHhyNWczXftPUliOjOGYQSUgd3gTe2ISeCmY1Sa
9wGrw2hiMz4tVLUs6C/zz9EV9F+RHQ2/96m0y6IPv0rAOWxtMTlvhG5FPrTEdHML
qz12/XcOpwKBgHhRH4HkGg2w6+kgCQoMSFrNFNBgEvGOVca+w6+G3S7DRZnU7BDB
bzXhY4zG02LVkkUEzmgf7u/y5tj0UdRTS3I2jRBJLVsjk4Po1NDxyWB7+S2kzvzV
LV1JY4z2LMdVO5O2Kunq41y9iywnXzCVxrClwEV9E/RfCBgQE7DG38RBAoGBAJDL
ARZmw98GaarJorUPE2V6AYnJ81Ao3jpfLO67ZlNa65rZUSa96MhbNV2j6zgSeTjk
Z1LTUG6wGZ9xuQ4lBriSgUNQppaVJiKj/J2EesuCLLCPDmsPKnqneMT7xTskISMT
pV4f9fp7hun7/cSwdahR3/laQDFe5x3swVZoqQybAoGAS1r1VrvBDABryleZrcbg
ul48VX9NwV5sBoi0BEqhBfUGnUQHRhGEMIX66U4vIOyd7QokiqikVIc+di1XcVLc
XHxT0g6qSJycfN0TbiRfxDhaC9h8P5ROJWhC02j68xSn2MtRjCEzJPgnBYHi/WlT
UpCbzjvpmE2CMpZTOBnfwSE=
-----END PRIVATE KEY-----`)
}

func newTestClientWithMock(httpClient goaliniex.HTTPClient) (*goaliniex.Client, error) {
	return goaliniex.NewClient(
		"https://sandbox.alixpay.com",
		"TEST_PARTNER",
		"TEST_SECRET",
		testPrivateKey(),
		goaliniex.WithHTTPClient(httpClient),
	)
}

func newTestClient(t *testing.T) *goaliniex.Client {
	t.Helper()

	privateKey, err := os.ReadFile("./alix-private-key.pem")
	if err != nil {
		t.Skipf("skipping test: unable to read private key: %v", err)
	}

	partnerCode := strings.TrimSpace(os.Getenv("ALIX_PARTNER_CODE"))
	if partnerCode == "" {
		t.Skip("skipping test: ALIX_PARTNER_CODE not set")
	}

	secretKey := strings.TrimSpace(os.Getenv("ALIX_SECRET_KEY"))
	if secretKey == "" {
		t.Skip("skipping test: ALIX_SECRET_KEY not set")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	client, err := goaliniex.NewClient(
		"https://sandbox.alixpay.com",
		partnerCode,
		secretKey,
		privateKey,
		goaliniex.WithDebug(true),
		goaliniex.WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}

func getTestEmail(t *testing.T) string {
	t.Helper()

	email := strings.TrimSpace(os.Getenv("ALIX_TEST_EMAIL"))
	if email == "" {
		t.Skip("skipping test: ALIX_TEST_EMAIL not set")
	}

	return email
}

func generateRandomGmail(t *testing.T) string {
	t.Helper()

	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		t.Fatalf("failed to generate random bytes: %v", err)
	}

	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("test.%s.%d@gmail.com", hex.EncodeToString(randomBytes), timestamp)
}

const testImageBase64 = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRof" +
	"Hh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIy" +
	"MjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAn" +
	"/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEB/wCwAB//2Q=="

func getTestImageDataURI() string {
	return "data:image/jpeg;base64," + testImageBase64
}

func generateTestSSN(t *testing.T) string {
	t.Helper()

	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		t.Fatalf("failed to generate random bytes: %v", err)
	}

	area := 900 + int(randomBytes[0])%100
	group := 1 + int(randomBytes[1])%99
	serial := 1 + int(randomBytes[2])%9999

	return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
}

func getTestEmail2(t *testing.T) string {
	t.Helper()

	email := strings.TrimSpace(os.Getenv("ALIX_TEST_EMAIL_2"))
	if email == "" {
		t.Skip("skipping test: ALIX_TEST_EMAIL_2 not set")
	}

	return email
}

func getTestEmails(t *testing.T) []string {
	t.Helper()

	email1 := strings.TrimSpace(os.Getenv("ALIX_TEST_EMAIL"))
	if email1 == "" {
		t.Skip("skipping test: ALIX_TEST_EMAIL not set")
	}

	emails := []string{email1}

	email2 := strings.TrimSpace(os.Getenv("ALIX_TEST_EMAIL_2"))
	if email2 != "" {
		emails = append(emails, email2)
	}

	return emails
}

func getTestBankCode(t *testing.T) string {
	t.Helper()

	bankCode := strings.TrimSpace(os.Getenv("ALIX_TEST_BANK_CODE"))
	if bankCode == "" {
		t.Skip("skipping test: ALIX_TEST_BANK_CODE not set")
	}

	return bankCode
}

func getTestBankAccountNumber(t *testing.T) string {
	t.Helper()

	bankAccountNumber := strings.TrimSpace(os.Getenv("ALIX_TEST_BANK_ACCOUNT_NUMBER"))
	if bankAccountNumber == "" {
		t.Skip("skipping test: ALIX_TEST_BANK_ACCOUNT_NUMBER not set")
	}

	return bankAccountNumber
}
