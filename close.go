package main

import (
	"fmt"
	"time"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Close(c *cli.Context) error {
	client := GetClient(c)

	item_ids := []string{}
	for _, arg := range c.Args().Slice() {
		item_ids = append(item_ids, arg)
	}

	if len(item_ids) == 0 {
		return CommandFailed
	}

	updatedItems := []todoist.Item{}
	closedCount := 0
	for _, item := range client.Store.Items {
		shouldClose := false
		for _, closeID := range item_ids {
			if item.ID == closeID {
				shouldClose = true
				closedCount++
				break
			}
		}
		if !shouldClose {
			updatedItems = append(updatedItems, item)
		}
	}
	client.Store.Items = updatedItems

	err := WriteCache(cachePath, client.Store)
	if err != nil {
		return fmt.Errorf("failed to update local cache: %w", err)
	}

	pc, err := GetPipelineCache(pipelineCachePath)
	if err != nil {
		return fmt.Errorf("failed to get pipeline cache: %w", err)
	}

	for _, itemID := range item_ids {
		command := todoist.NewCommand("item_close", map[string]interface{}{"id": itemID})

		pipelineItem := PipelineItem{
			IsClose:   true,
			CloseIDs:  []string{itemID},
			Command:   command,
			CreatedAt: time.Now(),
		}

		err = pc.AddItem(pipelineItem)
		if err != nil {
			return fmt.Errorf("failed to add close action to pipeline cache: %w", err)
		}
	}

	err = WritePipelineCache(pipelineCachePath, pc)
	if err != nil {
		return fmt.Errorf("failed to write pipeline cache: %w", err)
	}

	fmt.Printf("Closed %d item(s) (syncing in background)\n", closedCount)
	StartBackgroundSync(client, pipelineCachePath, cachePath)

	return nil
}
