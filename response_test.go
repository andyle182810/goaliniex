package goaliniex_test

import (
	"context"
	"testing"

	"github.com/andyle182810/goaliniex/v2"
)

func TestResponse_IsSuccess(t *testing.T) {
	t.Parallel()

	success := &goaliniex.Response[struct{}]{Code: goaliniex.CodeSuccess}
	if !success.IsSuccess() {
		t.Error("expected IsSuccess true for SUCCESS code")
	}

	failure := &goaliniex.Response[struct{}]{Code: goaliniex.CodeInvalidRequestData}
	if failure.IsSuccess() {
		t.Error("expected IsSuccess false for an error code")
	}
}

// Envelope decoding must not depend on the HTTP status: v3 says status is
// only a transport hint and callers branch on `code`. A 4xx with a
// well-formed envelope body must still decode successfully.
func TestEnvelopeDecodesRegardlessOfHTTPStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		statusCode int
		body       string
		wantCode   string
	}{
		{"200 success", 200, `{"code":"SUCCESS","message":"ok","data":{"version":"1.0.0","env":"sandbox"}}`, goaliniex.CodeSuccess},
		{
			"400 business error",
			400,
			`{"code":"INVALID_REQUEST_DATA","message":"bad","data":null}`,
			goaliniex.CodeInvalidRequestData,
		},
		{
			"500 internal error",
			500,
			`{"code":"INTERNAL_SERVER_ERROR","message":"boom","data":null}`,
			goaliniex.CodeInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockHTTPClient{response: mockResponse(testCase.statusCode, testCase.body)}

			client, err := newTestClientWithMock(mock)
			if err != nil {
				t.Fatalf("newTestClientWithMock: %v", err)
			}

			resp, err := client.GetPublicVersion(context.Background())
			if err != nil {
				t.Fatalf("GetPublicVersion: %v", err)
			}

			if resp.Code != testCase.wantCode {
				t.Errorf("Code = %q, want %q", resp.Code, testCase.wantCode)
			}
		})
	}
}
