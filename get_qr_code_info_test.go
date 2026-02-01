package goaliniex_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/andyle182810/goaliniex"
)

var errMockNetworkFailure = errors.New("network error")

func validateQRCodeBasicFields(t *testing.T, data *goaliniex.QRCodeInfo) {
	t.Helper()

	if data.BankAccountNumber != "888812345678" {
		t.Errorf("expected BankAccountNumber='888812345678', got %q", data.BankAccountNumber)
	}

	if data.BankCode != "Techcombank" {
		t.Errorf("expected BankCode='Techcombank', got %q", data.BankCode)
	}

	if data.BankName != "Ngân hàng TMCP Kỹ thương Việt Nam" {
		t.Errorf("expected BankName='Ngân hàng TMCP Kỹ thương Việt Nam', got %q", data.BankName)
	}

	if data.CountryCode != goaliniex.CountryCodeVN {
		t.Errorf("expected CountryCode=VN, got %v", data.CountryCode)
	}

	if data.QRType != goaliniex.QRTypeVietQR {
		t.Errorf("expected QRType=vietqr, got %v", data.QRType)
	}
}

func validateAdditionalData(t *testing.T, data map[string]any) {
	t.Helper()

	if data == nil {
		t.Fatal("AdditionalData is nil")
	}

	if len(data) != 2 {
		t.Errorf("expected 2 additional data fields, got %d", len(data))
	}

	bankAccNum, found := data["bankAccountNumber"].(string)
	if !found || bankAccNum != "888812345678" {
		t.Errorf("expected AdditionalData.bankAccountNumber='888812345678', got %v", data["bankAccountNumber"])
	}

	bankCode, exists := data["bankCode"].(string)
	if !exists || bankCode != "Techcombank" {
		t.Errorf("expected AdditionalData.bankCode='Techcombank', got %v", data["bankCode"])
	}
}

func TestClient_GetQRCodeInfo_Success(t *testing.T) {
	t.Parallel()

	successResponse := `{
		"success": true,
		"message": "Your request has been successful",
		"data": {
			"bankAccountNumber": "888812345678",
			"bankCode": "Techcombank",
			"bankName": "Ngân hàng TMCP Kỹ thương Việt Nam",
			"countryCode": "VN",
			"qrType": "vietqr",
			"additionalData": {
				"bankAccountNumber": "888812345678",
				"bankCode": "Techcombank"
			}
		},
		"errorCode": 0
	}`

	mockHTTP := &mockHTTPClient{
		response: mockResponse(http.StatusOK, successResponse), //nolint:bodyclose // mock response
		err:      nil,
	}

	client, err := newTestClientWithMock(mockHTTP)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if !resp.Success {
		t.Errorf("expected Success=true, got Success=%v", resp.Success)
	}

	if resp.Message != "Your request has been successful" {
		t.Errorf("expected Message='Your request has been successful', got Message=%q", resp.Message)
	}

	if resp.ErrorCode != 0 {
		t.Errorf("expected ErrorCode=0, got ErrorCode=%d", resp.ErrorCode)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	validateQRCodeBasicFields(t, resp.Data)
	validateAdditionalData(t, resp.Data.AdditionalData)
}

func TestClient_GetQRCodeInfo_UnsupportedQRCode(t *testing.T) {
	t.Parallel()

	failedResponse := `{
		"success": false,
		"message": "The QR code has not support yet.",
		"data": null,
		"errorCode": 33
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, failedResponse), //nolint:bodyclose // mock response with NopCloser
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "unsupported-qr-format",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if resp.Success {
		t.Errorf("expected Success=false, got Success=%v", resp.Success)
	}

	if resp.Message != "The QR code has not support yet." {
		t.Errorf("expected Message='The QR code has not support yet.', got Message=%q", resp.Message)
	}

	if resp.ErrorCode != 33 {
		t.Errorf("expected ErrorCode=33, got ErrorCode=%d", resp.ErrorCode)
	}

	if resp.Data != nil {
		t.Errorf("expected nil data for failed response, got %+v", resp.Data)
	}
}

func TestClient_GetQRCodeInfo_HTTPError(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: nil,
		err:      errMockNetworkFailure,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	_, err = client.GetQRCodeInfo(ctx, req)
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}

	if !errors.Is(err, goaliniex.ErrHTTPFailure) {
		t.Errorf("expected error to be ErrHTTPFailure, got %v", err)
	}
}

func TestClient_GetQRCodeInfo_InvalidJSON(t *testing.T) {
	t.Parallel()

	invalidJSONResponse := `{invalid json}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, invalidJSONResponse), //nolint:bodyclose // mock response with NopCloser
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	_, err = client.GetQRCodeInfo(ctx, req)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestClient_GetQRCodeInfo_EmptyQRContent(t *testing.T) {
	t.Parallel()

	failedResponse := `{
		"success": false,
		"message": "QR content is required",
		"data": null,
		"errorCode": 1
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, failedResponse), //nolint:bodyclose // mock response with NopCloser
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if resp.Success {
		t.Errorf("expected Success=false for empty QR content, got Success=%v", resp.Success)
	}
}

func TestClient_GetQRCodeInfo_ContextTimeout(t *testing.T) {
	t.Parallel()

	successResponse := `{
		"success": true,
		"message": "Your request has been successful",
		"data": {
			"bankAccountNumber": "888812345678",
			"bankCode": "Techcombank",
			"bankName": "Ngân hàng TMCP Kỹ thương Việt Nam",
			"countryCode": "VN",
			"qrType": "vietqr",
			"additionalData": {}
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, successResponse), //nolint:bodyclose // mock response with NopCloser
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected Success=true, got Success=%v", resp.Success)
	}
}

func TestClient_GetQRCodeInfo_HTTPStatusError(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusBadRequest, `{"error": "bad request"}`), //nolint:bodyclose // mock response
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "invalid",
	}

	_, err = client.GetQRCodeInfo(ctx, req)
	if err == nil {
		t.Fatal("expected error for non-200 status code, got nil")
	}

	if !errors.Is(err, goaliniex.ErrUnexpectedStatus) {
		t.Errorf("expected error to be ErrUnexpectedStatus, got %v", err)
	}
}

func TestClient_GetQRCodeInfo_ZeroAmount(t *testing.T) {
	t.Parallel()

	responseWithZeroAmount := `{
		"success": true,
		"message": "Your request has been successful",
		"data": {
			"bankAccountNumber": "888812345678",
			"bankCode": "Techcombank",
			"bankName": "Ngân hàng TMCP Kỹ thương Việt Nam",
			"countryCode": "VN",
			"qrType": "vietqr",
			"amount": 0,
			"additionalData": {}
		},
		"errorCode": 0
	}`

	mockClient := &mockHTTPClient{
		response: mockResponse(http.StatusOK, responseWithZeroAmount), //nolint:bodyclose // mock response with NopCloser
		err:      nil,
	}

	client, err := newTestClientWithMock(mockClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &goaliniex.GetQRCodeInfoRequest{
		QRContent: "00020101021138560010A0000007270126000697040701128888123456780208QRIBFTTA53037045802VN63042249",
	}

	resp, err := client.GetQRCodeInfo(ctx, req)
	if err != nil {
		t.Fatalf("GetQRCodeInfo returned error: %v", err)
	}

	if resp.Data.Amount != 0 {
		t.Errorf("expected Amount=0, got Amount=%f", resp.Data.Amount)
	}
}

func TestClient_GetQRCodeInfo_DifferentQRTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		qrType     goaliniex.QRType
		jsonQRType string
	}{
		{
			name:       "VietQR",
			qrType:     goaliniex.QRTypeVietQR,
			jsonQRType: "vietqr",
		},
		{
			name:       "Philippine PPMI P2M",
			qrType:     goaliniex.QRTypePHPPMIP2M,
			jsonQRType: "ph.ppmi.p2m",
		},
		{
			name:       "PIX",
			qrType:     goaliniex.QRTypePIX,
			jsonQRType: "pix",
		},
		{
			name:       "PayWithCrypto",
			qrType:     goaliniex.QRTypePayWithCrypto,
			jsonQRType: "paywithcrypto",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			response := `{
				"success": true,
				"message": "Your request has been successful",
				"data": {
					"bankAccountNumber": "123456789",
					"bankCode": "TEST",
					"bankName": "Test Bank",
					"countryCode": "VN",
					"qrType": "` + testCase.jsonQRType + `",
					"additionalData": {}
				},
				"errorCode": 0
			}`

			mockClient := &mockHTTPClient{
				response: mockResponse(http.StatusOK, response), //nolint:bodyclose // mock response with NopCloser
				err:      nil,
			}

			client, err := newTestClientWithMock(mockClient)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &goaliniex.GetQRCodeInfoRequest{
				QRContent: "test-qr-content",
			}

			resp, err := client.GetQRCodeInfo(ctx, req)
			if err != nil {
				t.Fatalf("GetQRCodeInfo returned error: %v", err)
			}

			if resp.Data.QRType != testCase.qrType {
				t.Errorf("expected QRType=%v, got %v", testCase.qrType, resp.Data.QRType)
			}
		})
	}
}
