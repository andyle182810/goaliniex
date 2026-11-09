package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type WidgetThemeResponse struct {
	Theme map[string]any `json:"theme"`
}

// GetWidgetTheme returns the caller partner's current alix-ramp-widget
// theme override, a sparse JSON merge patch to apply onto the widget
// default. Requires a bearer token (not public), GET /v3/widget/theme.
func (c *Client) GetWidgetTheme(ctx context.Context, userEmail string) (*Response[WidgetThemeResponse], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/widget/theme",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[WidgetThemeResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
