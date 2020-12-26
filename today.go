package main

import (
	"flag"

	"github.com/urfave/cli"
)

func Today(c *cli.Context) error {
	flagSet := flag.NewFlagSet(c.Command.Name, flag.ContinueOnError)
	_ = flagSet.String("filter", "", "filter")

	listCtx := cli.NewContext(c.App, flagSet, c)
	err := listCtx.Set("filter", "today")
	if err != nil {
		// TODO: send error to system log or something?
	}

	return List(listCtx)
}
