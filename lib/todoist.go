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
	AccessToken    string
	DebugMode      bool
	Color          bool
	DateFormat     string
	DateTimeFormat string
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

	c.Log("config: %#v", c.config)

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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req = req.WithContext(ctx)

	c.Log("request: %#v", req)
	c.Log("request.URL: %#v", req.URL)
	c.Log("params: %#v", body)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Log("response: %#v", resp)

	if resp.StatusCode != http.StatusOK {
		c.Log(ParseAPIError("bad request", resp).Error())
		return ParseAPIError("bad request", resp)
	} else if res == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&res)
}

func (c *Client) doRestApi(ctx context.Context, method string, uri string, body interface{}, res interface{}) error {
	c.Log("doRestApi: called")
	u, err := url.Parse(Server)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, uri)

	c.Log("config: %#v", c.config)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req = req.WithContext(ctx)

	c.Log("request: %#v", req)
	c.Log("request.URL: %#v", req.URL)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Log("response: %#v", resp)

	if resp.StatusCode != http.StatusOK {
		c.Log(ParseAPIError("bad request", resp).Error())
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

func (c *Client) QuickCommand(ctx context.Context, text string) error {
	var item Item

	body := map[string]interface{}{
		"text": text,
	}

	return c.doRestApi(ctx, http.MethodPost, "tasks/quick", body, &item)
}

func (c *Client) Sync(ctx context.Context) error {
	params := url.Values{"sync_token": {"*"}, "resource_types": {"[\"all\"]"}}

	err := c.doApi(ctx, http.MethodPost, "sync", params, &c.Store)
	if err != nil {
		return err
	}
	c.Store.ConstructItemTree()
	return nil
}

func (c *Client) CompleteItemIDByPrefix(prefix string) (id string, err error) {
	var matchid string = ""
	for _, cmpid := range c.Store.Items {
		if strings.HasPrefix(cmpid.GetID(), prefix) {
			if matchid != "" {
				// Ambiguous prefix, return converted input instead
				return prefix, nil
			} else {
				matchid = cmpid.GetID()
			}
		}
	}
	if matchid != "" {
		return matchid, nil
	} else {
		return prefix, nil
	}
}

// doRestApiV2 makes requests to the Todoist REST API v2
// Used for endpoints like /reminders that are not part of the Sync API
func (c *Client) doRestApiV2(ctx context.Context, method string, uri string, body interface{}, res interface{}) error {
	c.Log("doRestApiV2: called")
	u, err := url.Parse(RestV2Server)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, uri)

	c.Log("config: %#v", c.config)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req = req.WithContext(ctx)

	c.Log("request: %#v", req)
	c.Log("request.URL: %#v", req.URL)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Log("response: %#v", resp)

	// Accept both 200 OK and 204 No Content as success
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		c.Log(ParseAPIError("bad request", resp).Error())
		return ParseAPIError("bad request", resp)
	}

	if res != nil && resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(res)
	}
	return nil
}
