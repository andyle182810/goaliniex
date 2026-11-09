package goaliniex

import (
	"context"
	"encoding/json"
	"net/http"
)

type WidgetTranslationsResponse struct {
	Translations map[string]any `json:"translations"`
}

// GetWidgetTranslations returns the caller partner's current
// alix-ramp-widget translations override, a sparse JSON merge patch keyed
// by namespace. Requires a bearer token (not public),
// GET /v3/widget/translations.
func (c *Client) GetWidgetTranslations(
	ctx context.Context,
	userEmail string,
) (*Response[WidgetTranslationsResponse], error) {
	apiRequest := request{
		Method:    http.MethodGet,
		Endpoint:  "/v3/widget/translations",
		UserEmail: userEmail,
	}

	rawResponse, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	response := new(Response[WidgetTranslationsResponse])
	if err := json.Unmarshal(rawResponse, response); err != nil {
		return nil, err
	}

	return response, nil
}
