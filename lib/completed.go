package todoist

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type Completed struct {
	Items    CompletedItems `json:"items"`
	Projects interface{}    `json:"projects"`
}

func (c *Client) CompletedAll(ctx context.Context, r *Completed) error {
	// v1 API requires since/until parameters in UTC with "Z" suffix (max 3 month range)
	// Default to last 30 days
	now := time.Now().UTC()
	since := now.AddDate(0, 0, -30).Format("2006-01-02T15:04:05Z")
	until := now.Format("2006-01-02T15:04:05Z")

	params := url.Values{
		"since": {since},
		"until": {until},
	}

	return c.doApi(ctx, http.MethodGet, "tasks/completed/by_completion_date", params, &r)
}
