package main

import (
	"context"

	"github.com/urfave/cli"
)

func Quick(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	client.QuickCommand(context.Background(), c.Args().First())
	return nil
}
