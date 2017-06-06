package main

import (
	"context"
	"strconv"

	"github.com/urfave/cli"
)

func Delete(c *cli.Context) error {
	client := GetClient(c)

	item_ids := []int{}
	for _, arg := range c.Args() {
		item_id, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		item_ids = append(item_ids, item_id)
	}

	if len(item_ids) == 0 {
		return CommandFailed
	}

	if err := client.DeleteItem(context.Background(), item_ids); err != nil {
		return CommandFailed
	}

	if err := Sync(c); err != nil {
		return CommandFailed
	}

	return nil
}
