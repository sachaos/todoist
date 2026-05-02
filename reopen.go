package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func Reopen(c *cli.Context) error {
	client := GetClient(c)

	for _, id := range c.Args().Slice() {
		if err := client.ReopenItem(context.Background(), id); err != nil {
			return fmt.Errorf("failed to reopen task %s: %w", id, err)
		}
	}

	if c.Args().Len() == 0 {
		return fmt.Errorf("no task IDs provided\nUsage: todoist reopen <Item ID> [<Item ID>...]\nUse `todoist completed-list` to find IDs of recently closed tasks")
	}

	// No Sync(c): completed tasks are never in the local cache, so the cache
	// is not stale after a reopen. Run `todoist sync` to see the task in `list`.
	return nil
}
