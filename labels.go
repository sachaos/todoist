package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func Labels(sync todoist.Sync, c *cli.Context) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	defer writer.Flush()

	for _, label := range sync.Labels {
		writer.Write([]string{IdFormat(label), "@" + label.Name})
	}

	return nil
}
