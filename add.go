package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Add(config Config, c *cli.Context) error {
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

	err = Sync(config, c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
