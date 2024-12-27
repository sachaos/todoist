package todoist

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/urfave/cli/v2"
)

type Completed struct {
	Items    CompletedItems `json:"items"`
	Projects interface{}    `json:"projects"`
}

func (c *Client) CompletedAll(cli *cli.Context, ctx context.Context, r *Completed) error {
	v := url.Values{}

	v.Add("limit", strconv.Itoa(cli.Int("limit")))

	if since := cli.String("since"); since != "" {
		v.Add("since", cli.String("since"))
	}

	return c.doApi(ctx, http.MethodPost, "completed/get_all", v, &r)
}
