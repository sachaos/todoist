package main

import (
	todoist "github.com/sachaos/todoist/lib"
)

func TrySyncPendingItems(client *todoist.Client) {
	go BackgroundSyncWorker(client, pipelineCachePath, cachePath)
}
