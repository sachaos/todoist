package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
)

func Close(config Config, c *cli.Context) error {
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

	err := lib.CloseItem(item_ids, config.Token)
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(config, c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
