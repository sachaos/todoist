package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
)

func Close(config Config, c *cli.Context) error {
	var sync lib.Sync
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

	sync, err = lib.FetchCache(config.Token)
	if err != nil {
		return CommandFailed
	}
	err = lib.SaveCache(default_cache_path, sync)
	if err != nil {
		return CommandFailed
	}

	return nil
}
