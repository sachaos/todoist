package main

import (
	"context"

	"github.com/urfave/cli/v2"
)

func DeleteSection(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	sectionID := c.Args().First()
	if client.Store.FindSection(sectionID) == nil {
		return IdNotFound
	}

	if err := client.DeleteSection(context.Background(), sectionID); err != nil {
		return err
	}

	return Sync(c)
}
