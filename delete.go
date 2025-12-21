package main

import (
	"fmt"
	"time"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Delete(c *cli.Context) error {
	client := GetClient(c)

	item_ids := []string{}
	for _, arg := range c.Args().Slice() {
		item_id, err := client.CompleteItemIDByPrefix(arg)
		if err != nil {
			return err
		}
		item_ids = append(item_ids, item_id)
	}

	if len(item_ids) == 0 {
		return CommandFailed
	}

	// Update local cache immediately by removing deleted items
	updatedItems := []todoist.Item{}
	deletedCount := 0
	for _, item := range client.Store.Items {
		shouldDelete := false
		for _, deleteID := range item_ids {
			if item.ID == deleteID {
				shouldDelete = true
				deletedCount++
				break
			}
		}
		if !shouldDelete {
			updatedItems = append(updatedItems, item)
		}
	}
	client.Store.Items = updatedItems

	err := WriteCache(cachePath, client.Store)
	if err != nil {
		return fmt.Errorf("failed to update local cache: %w", err)
	}

	// Add delete commands to pipeline cache
	pc, err := GetPipelineCache(pipelineCachePath)
	if err != nil {
		return fmt.Errorf("failed to get pipeline cache: %w", err)
	}

	for _, itemID := range item_ids {
		command := todoist.NewCommand("item_delete", map[string]interface{}{"id": itemID})

		pipelineItem := PipelineItem{
			Command:   command,
			CreatedAt: time.Now(),
		}

		err = pc.AddItem(pipelineItem)
		if err != nil {
			return fmt.Errorf("failed to add delete action to pipeline cache: %w", err)
		}
	}

	err = WritePipelineCache(pipelineCachePath, pc)
	if err != nil {
		return fmt.Errorf("failed to write pipeline cache: %w", err)
	}

	fmt.Printf("Deleted %d item(s) (syncing in background)\n", deletedCount)
	StartBackgroundSync(client, pipelineCachePath, cachePath)

	return nil
}
