package lib

import (
	"errors"
)

var (
	PostFailed = errors.New("Post Failed")
	FindFailed = errors.New("Find Failed")
	SyncFailed = errors.New("Sync Failed")
)
