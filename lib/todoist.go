package todoist

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Config struct {
	AccessToken string
	DebugMode   bool
}

type Client struct {
	http.Client
	config *Config
	Store  *Store
}

func NewClient(config *Config) *Client {
	return &Client{
		Client: *http.DefaultClient,
		config: config,
	}
}

func (c *Client) Log(format string, v ...interface{}) {
	if c.config.DebugMode {
		log.Printf(format, v...)
	}
}

func (c *Client) doApi(ctx context.Context, method string, uri string, params url.Values, res interface{}) error {
	c.Log("doAPi: called")
	u, err := url.Parse(Server)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, uri)

	params.Add("token", c.config.AccessToken)

	var body io.Reader
	if method == http.MethodGet {
		u.RawQuery = params.Encode()
	} else {
		body = strings.NewReader(params.Encode())
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ParseAPIError("bad request", resp)
	} else if res == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&res)
}

type ExecResult struct {
	SyncToken     string      `json:"sync_token"`
	SyncStatus    interface{} `json:"sync_status"`
	TempIdMapping interface{} `json:"temp_id_mapping"`
}

func (c *Client) ExecCommands(ctx context.Context, commands Commands) error {
	var r ExecResult
	return c.doApi(ctx, http.MethodPost, "sync", commands.UrlValues(), &r)
}

func (c *Client) Sync(ctx context.Context) error {
	params := url.Values{"sync_token": {"*"}, "resource_types": {"[\"all\"]"}}

	err := c.doApi(ctx, http.MethodPost, "sync", params, &c.Store)
	if err != nil {
		return err
	}
	c.Store.ConstructItemOrder()
	return nil
}
