package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// QRType is the detected QR scheme, named {country}_{scheme}. Full v3 value
// set, replacing v2's ad-hoc lowercase strings.
type QRType string

const (
	QRTypeVNVietQR    QRType = "VN_VIETQR"
	QRTypeVNVNPayQR   QRType = "VN_VNPAY_QR"
	QRTypePHQRPH      QRType = "PH_QRPH"
	QRTypeIDQRIS      QRType = "ID_QRIS"
	QRTypeBRPix       QRType = "BR_PIX"
	QRTypeARQR3       QRType = "AR_QR3"
	QRTypePEQR        QRType = "PE_QR"
	QRTypeTHPromptPay QRType = "TH_PROMPTPAY"
	QRTypeMYDuitNow   QRType = "MY_DUITNOW"
	QRTypeSGPayNow    QRType = "SG_PAYNOW"
	QRTypeKHKHQR      QRType = "KH_KHQR"
	QRTypeINUPI       QRType = "IN_UPI"
	QRTypeHKFPS       QRType = "HK_FPS"
	QRTypePYSIPQR     QRType = "PY_SIP_QR"
	QRTypeNGNQR       QRType = "NG_NQR"
	QRTypeMXCodi      QRType = "MX_CODI"
	QRTypeUnknown     QRType = "UNKNOWN"
)

// CountryCode is the country of the decoded QR scheme. This is a smaller
// set than the fiat-currency country list — QR decoding only covers markets
// with a local QR scheme today.
type CountryCode string

const (
	CountryCodeVN      CountryCode = "VN"
	CountryCodePH      CountryCode = "PH"
	CountryCodeID      CountryCode = "ID"
	CountryCodeAR      CountryCode = "AR"
	CountryCodeBR      CountryCode = "BR"
	CountryCodePE      CountryCode = "PE"
	CountryCodeUnknown CountryCode = "UNKNOWN"
)

// QRPayeeAddressing describes how the recipient is identified inside a
// decoded QR code.
type QRPayeeAddressing string

const (
	QRPayeeAddressingAccountIdentifier QRPayeeAddressing = "ACCOUNT_IDENTIFIER"
	QRPayeeAddressingProxy             QRPayeeAddressing = "PROXY"
	QRPayeeAddressingSchemeToken       QRPayeeAddressing = "SCHEME_TOKEN"
)

// QRTransferType distinguishes merchant vs person QR payments.
type QRTransferType string

const (
	QRTransferTypeP2M QRTransferType = "P2M"
	QRTransferTypeP2P QRTransferType = "P2P"
)

type DecodeQRRequest struct {
	QR string `json:"qr"`
}

type DecodeQRResponse struct {
	Supported       bool              `json:"supported"`
	PayeeAddressing QRPayeeAddressing `json:"payeeAddressing"`
	BankCode        string            `json:"bankCode"`
	BankName        string            `json:"bankName"`
	AccountNumber   string            `json:"accountNumber"`
	AccountName     string            `json:"accountName"`
	Description     string            `json:"description"`
	Amount          float64           `json:"amount"`
	// QRToken is a signed opaque token carrying the decoded payment data;
	// pass it back verbatim as ScanToPayOrderRequest.QRToken.
	QRToken      string         `json:"qrToken"`
	FiatCurrency FiatCurrency   `json:"fiatCurrency"`
	QRType       QRType         `json:"qrType"`
	TransferType QRTransferType `json:"transferType"`
	CountryCode  CountryCode    `json:"countryCode"`
}

// Error codes specific to QR decoding.
const (
	CodeQRCodeNotSupported            = "QR_CODE_NOT_SUPPORTED"
	CodeServiceTemporarilyUnavailable = "SERVICE_TEMPORARILY_UNAVAILABLE"
)

// DecodeQR decodes a payment QR code into its bank and account fields. This
// is an open (public, no bearer token) endpoint, POST
// /v3/scan-to-pay/qr-decodings.
func (c *Client) DecodeQR(ctx context.Context, req *DecodeQRRequest) (*Response[DecodeQRResponse], error) {
	apiRequest := request{
		Method:   http.MethodPost,
		Endpoint: "/v3/scan-to-pay/qr-decodings",
		Params:   req,
		Public:   true,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[DecodeQRResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
