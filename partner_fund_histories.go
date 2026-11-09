package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// FundHistoryEntry is one entry in the partner's crypto fund ledger. Amount
// is signed: negative when spent, positive when topped up or refunded.
type FundHistoryEntry struct {
	ID           int     `json:"id"`
	Type         string  `json:"type"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balanceAfter"`
	CreatedAt    string  `json:"createdAt"`
}

// GetPartnerFundHistoriesRequest fields are all optional. Limit/Offset
// paginate the ledger; there is no response-side pagination wrapper.
type GetPartnerFundHistoriesRequest struct {
	Currency string `json:"currency,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// GetPartnerFundHistories lists entries from the partner's crypto fund
// ledger, GET /v3/partner/fund/histories.
func (c *Client) GetPartnerFundHistories(
	ctx context.Context,
	req *GetPartnerFundHistoriesRequest,
) (*Response[[]FundHistoryEntry], error) {
	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/v3/partner/fund/histories",
		Params:   req,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]FundHistoryEntry])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
