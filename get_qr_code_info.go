package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type QRType string

const (
	QRTypeVietQR        QRType = "vietqr"
	QRTypePHPPMIP2M     QRType = "ph.ppmi.p2m"
	QRTypeComP2PQRPay   QRType = "com.p2pqrpay"
	QRTypePIX           QRType = "pix"
	QRTypeQR3           QRType = "qr3"
	QRTypePayWithCrypto QRType = "paywithcrypto"
)

type CountryCode string

const (
	CountryCodeVN CountryCode = "VN"
	CountryCodePH CountryCode = "PH"
	CountryCodeTH CountryCode = "TH"
	CountryCodeGE CountryCode = "GE"
	CountryCodeBR CountryCode = "BR"
	CountryCodeAR CountryCode = "AR"
	CountryCodePE CountryCode = "PE"
)

type GetQRCodeInfoRequest struct {
	QRContent string `json:"qrContent"`
}

type QRCodeInfo struct {
	BankAccountNumber string         `json:"bankAccountNumber"`
	BankCode          string         `json:"bankCode"`
	BankName          string         `json:"bankName"`
	CountryCode       CountryCode    `json:"countryCode"`
	QRType            QRType         `json:"qrType"`
	AdditionalData    map[string]any `json:"additionalData"`
	Amount            float64        `json:"amount"`
}

func (c *Client) GetQRCodeInfo(ctx context.Context, req *GetQRCodeInfoRequest) (*Response[QRCodeInfo], error) {
	apiRequest := request{
		Method:      http.MethodGet,
		Endpoint:    "/api/v2/public/get-qr-code-info",
		Params:      req,
		SigningData: nil,
		Header:      nil,
		Body:        nil,
		FullURL:     "",
		Public:      true,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[QRCodeInfo])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
