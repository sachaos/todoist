package main

import (
	"context"

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

	if err := client.CloseItem(context.Background(), item_ids); err != nil {
		return err
	}

	return Sync(c)
}
