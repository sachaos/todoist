package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	FindFailed = errors.New("Find Failed")
)

const (
	Server       = "https://api.todoist.com/api/v1/"
	RestV2Server = "https://api.todoist.com/rest/v2/"
)

func ParseAPIError(prefix string, resp *http.Response) error {
	errMsg := fmt.Sprintf("%s: %s", prefix, resp.Status)
	var e struct {
		Error string `json:"error"`
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	json.Unmarshal(bodyBytes, &e)
	if e.Error != "" {
		errMsg = fmt.Sprintf("%s: %s", errMsg, e.Error)
	} else if len(bodyBytes) > 0 {
		errMsg = fmt.Sprintf("%s: %s", errMsg, string(bodyBytes))
	}

	return errors.New(errMsg)
}
