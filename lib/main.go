package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	FindFailed = errors.New("Find Failed")
)

const (
	Server     = "https://todoist.com/API/v8/"
	DateFormat = "Mon 2 Jan 2006 15:04:05 +0000"
)

func ParseAPIError(prefix string, resp *http.Response) error {
	errMsg := fmt.Sprintf("%s: %s", prefix, resp.Status)
	var e struct {
		Error string `json:"error"`
	}

	json.NewDecoder(resp.Body).Decode(&e)
	if e.Error != "" {
		errMsg = fmt.Sprintf("%s: %s", errMsg, e.Error)
	}

	return errors.New(errMsg)
}
