package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
)

func Modify(config Config, c *cli.Context) error {
	item := lib.Item{}

	if !c.Args().Present() {
		return CommandFailed
	}

	item_id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return err
	}
	item.ID = item_id
	item.Priority = c.Int("priority")
	item.Content = c.String("content")
	err = lib.UpdateItem(item, config.Token)
	if err != nil {
		return CommandFailed
	}

	err = Sync(config, c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
