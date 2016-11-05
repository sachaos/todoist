package lib

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	FindFailed = errors.New("Find Failed")
)

const (
	DateFormat = "Mon 2 Jan 2006 15:04:05 +0000"
)

func SyncRequest(params url.Values) ([]byte, error) {
	resp, err := http.PostForm("https://todoist.com/API/v7/sync", params)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return body, nil
}
