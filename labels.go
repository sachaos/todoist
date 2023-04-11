package main

import (
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func Labels(c *cli.Context) error {
	client := GetClient(c)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	defer writer.Flush()

	if c.Bool("header") {
		writer.Write([]string{"ID", "Name"})
	}

	for _, label := range client.Store.Labels {
		writer.Write([]string{IdFormat(label), "@" + label.Name})
	}

	return nil
}
