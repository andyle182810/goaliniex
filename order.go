package goaliniex

// OrderStatus is an order's lifecycle status, shared across scan-to-pay,
// cashout, off-ramp and on-ramp orders. Spelled PROCESSING (not PROGRESSING,
// which is a distinct KycStatus value).
type OrderStatus string

const (
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusSuccess    OrderStatus = "SUCCESS"
	OrderStatusFailed     OrderStatus = "FAILED"
	OrderStatusStuck      OrderStatus = "STUCK"
)

// PayoutKind identifies how a fiat payout recipient is addressed.
type PayoutKind string

const (
	PayoutKindBankAccount PayoutKind = "BANK_ACCOUNT"
	PayoutKindQRPayment   PayoutKind = "QR_PAYMENT"
)

// PayoutBankAccount and PayoutQRRecipient sub-field names are inferred from
// the equivalent v2 BankTransfer shape; the v3 doc declares bankAccount/
// qrRecipient as opaque required objects without enumerating their fields.
// Confirm against the live openapi/v3.json before relying on exact names.
type PayoutBankAccount struct {
	BankCode      string `json:"bankCode"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

type PayoutQRRecipient struct {
	QRType        string `json:"qrType"`
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

// Payout describes the fiat leg paid to the recipient of a scan-to-pay or
// cashout order.
type Payout struct {
	Kind        PayoutKind         `json:"kind"`
	BankAccount *PayoutBankAccount `json:"bankAccount,omitempty"`
	QRRecipient *PayoutQRRecipient `json:"qrRecipient,omitempty"`
	Description string             `json:"description"`
	Amount      float64            `json:"amount"`
}

// Quote is the locked-quote echo returned on an order created with a
// quoteId, shared across all four order products.
type Quote struct {
	QuoteID   string  `json:"quoteId"`
	Price     float64 `json:"price"`
	ExpiresAt string  `json:"expiresAt"`
}

// Fees is the processing/system fee breakdown used by off-ramp and on-ramp.
// scan-to-pay and cashout use a single top-level "fee" field instead.
type Fees struct {
	ProcessingFee float64 `json:"processingFee"`
	SystemFee     float64 `json:"systemFee"`
}

// Bank identifies a payout destination bank, shared by the scan-to-pay and
// cashout banks endpoints.
type Bank struct {
	BankCode string `json:"bankCode"`
	BankName string `json:"bankName"`
}

// Currency is the crypto debited/credited by scan-to-pay and cashout orders;
// both products only support USDT today. off-ramp/on-ramp's cryptoCurrency
// is open-ended and modeled as plain string instead.
type Currency string

const CurrencyUSDT Currency = "USDT"

// FiatCurrency is the fiat currency accepted across all v3 products.
type FiatCurrency string

const (
	FiatCurrencyUSD FiatCurrency = "USD"
	FiatCurrencyVND FiatCurrency = "VND"
	FiatCurrencyPHP FiatCurrency = "PHP"
	FiatCurrencyTHB FiatCurrency = "THB"
	FiatCurrencyIDR FiatCurrency = "IDR"
	FiatCurrencyBRL FiatCurrency = "BRL"
	FiatCurrencyGEL FiatCurrency = "GEL"
	FiatCurrencyARS FiatCurrency = "ARS"
	FiatCurrencyPEN FiatCurrency = "PEN"
	FiatCurrencyNGN FiatCurrency = "NGN"
	FiatCurrencyKES FiatCurrency = "KES"
	FiatCurrencyGHS FiatCurrency = "GHS"
	FiatCurrencyZAR FiatCurrency = "ZAR"
)
