package todoist

import (
	"encoding/json"
	"net/url"
)

type Completed struct {
	Items    CompletedItems `json:"items"`
	Projects interface{}    `json:"projects"`
}

func CompletedAll(token string) (Completed, error) {
	var completed Completed
	body, err := APIRequest("completed/get_all", url.Values{"token": {token}})
	err = json.Unmarshal(body, &completed)
	if err != nil {
		return Completed{}, err
	}
	return completed, nil
}
