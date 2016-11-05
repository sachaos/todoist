package lib

import (
	"errors"
)

var (
	PostFailed = errors.New("Post Failed")
	FindFailed = errors.New("Find Failed")
	SyncFailed = errors.New("Sync Failed")
)

const (
	DateFormat = "Mon 2 Jan 2006 15:04:05 +0000"
)
