package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	commit string
	date   string
)

func Version(c *cli.Context) error {
	if version == "" {
		fmt.Println("todoist (dev build)")
		return nil
	}
	if commit != "" && date != "" {
		fmt.Printf("todoist %s (commit %s, built %s)\n", version, commit, date)
		return nil
	}
	fmt.Printf("todoist %s\n", version)
	return nil
}
