package main

import (
	"context"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func AddProject(c *cli.Context) error {
	client := GetClient(c)

	project := todoist.Project{}
	if !c.Args().Present() {
		return CommandFailed
	}

	project.Name = c.Args().First()
	project.Color = c.String("color")
	project.ItemOrder = c.Int("item-order")

	if err := client.AddProject(context.Background(), project); err != nil {
		return err
	}

	return Sync(c)
}
