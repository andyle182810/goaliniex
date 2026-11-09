package goaliniex

// Response is the v3 response envelope: every response, success or error, is
// exactly {code, message, data}. Branch on Code, not on the HTTP status.
type Response[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

// IsSuccess reports whether the response's code is the success code.
func (r *Response[T]) IsSuccess() bool {
	return r.Code == CodeSuccess
}

// Common error codes, returned by every authenticated endpoint. Per-endpoint
// error codes are declared alongside the endpoint file that returns them.
const (
	CodeSuccess             = "SUCCESS"
	CodeInvalidRequestData  = "INVALID_REQUEST_DATA"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeAccessTokenInvalid  = "ACCESS_TOKEN_INVALID"

	CodePartnerNotFound       = "PARTNER_NOT_FOUND"
	CodePartnerAccountLocked  = "PARTNER_ACCOUNT_LOCKED"
	CodePartnerIPNotAllowed   = "PARTNER_IP_NOT_ALLOWED"
	CodeServiceNotAvailable   = "SERVICE_NOT_AVAILABLE"
	CodeServiceUnderMaint     = "SERVICE_UNDER_MAINTENANCE"
	CodeTransactionNotFound   = "TRANSACTION_NOT_FOUND"
	CodeDuplicateTransaction  = "DUPLICATE_TRANSACTION_ID"
	CodeQuoteNotFound         = "QUOTE_NOT_FOUND"
	CodeQuoteExpired          = "QUOTE_EXPIRED"
	CodeQuoteAlreadyUsed      = "QUOTE_ALREADY_USED"
	CodeQuoteMismatch         = "QUOTE_MISMATCH"
	CodeTxnAmountLimitExceed  = "TRANSACTION_AMOUNT_LIMIT_EXCEEDED"
	CodeDailyAmountExceed     = "DAILY_AMOUNT_LIMIT_EXCEEDED"
	CodeWeeklyAmountExceed    = "WEEKLY_AMOUNT_LIMIT_EXCEEDED"
	CodeMonthlyAmountExceed   = "MONTHLY_AMOUNT_LIMIT_EXCEEDED"
	CodePartnerInsufficientBl = "PARTNER_INSUFFICIENT_BALANCE"
)
