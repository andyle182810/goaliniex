package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

// DestinationFieldOption is a selectable choice for a "select"-type
// DestinationField.
type DestinationFieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// DestinationField describes one field the caller must collect to identify
// an off-ramp fiat recipient; Name matches the key in OffRampRecipient.
type DestinationField struct {
	Name    string                   `json:"name"`
	Type    string                   `json:"type"`
	Label   string                   `json:"label"`
	Options []DestinationFieldOption `json:"options,omitempty"`
}

// GetOffRampDestinationFields lists the fields required to identify the
// fiat recipient for an off-ramp order (replaces a plain bank list),
// GET /v3/off-ramp/destination-fields.
func (c *Client) GetOffRampDestinationFields(
	ctx context.Context,
	userEmail string,
	fiatCurrency FiatCurrency,
) (*Response[[]DestinationField], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/off-ramp/destination-fields",
		Params:    map[string]any{"fiatCurrency": fiatCurrency},
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]DestinationField])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
