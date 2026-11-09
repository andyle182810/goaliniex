package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// CreateWidgetLinkRequest fields are all optional. CryptoAmount and
// FiatAmount are strings (not numbers) per the spec's pre-fill semantics:
// invalid values are silently ignored by the widget rather than rejected.
type CreateWidgetLinkRequest struct {
	Type           string `json:"type,omitempty"`
	CryptoCurrency string `json:"cryptoCurrency,omitempty"`
	Network        string `json:"network,omitempty"`
	CryptoAmount   string `json:"cryptoAmount,omitempty"`
	FiatAmount     string `json:"fiatAmount,omitempty"`
	FiatCurrency   string `json:"fiatCurrency,omitempty"`
	Theme          string `json:"theme,omitempty"`
}

type WidgetLink struct {
	URL       string `json:"url"`
	ExpiresAt int64  `json:"expiresAt"`
}

// CreateWidgetLink creates a time-limited embeddable widget link for the
// user of the caller's access token, POST /v3/widget-links.
func (c *Client) CreateWidgetLink(
	ctx context.Context,
	userEmail string,
	req *CreateWidgetLinkRequest,
) (*Response[WidgetLink], error) {
	apiRequest := request{
		Method:    http.MethodPost,
		Endpoint:  "/v3/widget-links",
		Params:    req,
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[WidgetLink])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
