package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func UpdateSection(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	sectionID := c.Args().First()
	if client.Store.FindSection(sectionID) == nil {
		return IdNotFound
	}

	name := c.String("name")
	if name == "" {
		return fmt.Errorf("--name flag is required")
	}

	if err := client.UpdateSection(context.Background(), sectionID, name); err != nil {
		return err
	}

	return Sync(c)
}
