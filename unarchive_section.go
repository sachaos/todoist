package main

import (
	"context"

	"github.com/urfave/cli/v2"
)

func UnarchiveSection(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	sectionID := c.Args().First()

	// Don't validate against local cache - archived sections aren't cached locally
	if err := client.UnarchiveSection(context.Background(), sectionID); err != nil {
		return err
	}

	return Sync(c)
}
