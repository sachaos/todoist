package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Add(config Config, c *cli.Context) error {
	var sync lib.Sync
	item := lib.Item{}

	if !c.Args().Present() {
		return CommandFailed
	}

	item.Content = c.Args().First()
	item.Priority = 1
	err := lib.AddItem(item, config.Token)
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
