package main

import (
	"context"

	"github.com/urfave/cli/v2"
)

func Reopen(c *cli.Context) error {
	client := GetClient(c)

	for _, id := range c.Args().Slice() {
		if err := client.ReopenItem(context.Background(), id); err != nil {
			return err
		}
	}

	if c.Args().Len() == 0 {
		return CommandFailed
	}

	// No Sync(c): completed tasks are never in the local cache, so the cache
	// is not stale after a reopen. Run `todoist sync` to see the task in `list`.
	return nil
}
