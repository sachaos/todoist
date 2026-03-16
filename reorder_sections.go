package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func ReorderSections(c *cli.Context) error {
	client := GetClient(c)

	if c.Args().Len() < 2 {
		return fmt.Errorf("reorder-sections requires at least 2 section IDs")
	}

	ids := c.Args().Slice()
	for _, id := range ids {
		if client.Store.FindSection(id) == nil {
			return fmt.Errorf("section id not found: %s", id)
		}
	}

	if err := client.ReorderSections(context.Background(), ids); err != nil {
		return err
	}

	return Sync(c)
}
