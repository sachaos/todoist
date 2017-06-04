package main

import (
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

func Labels(c *cli.Context) error {
	client := GetClient(c)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	defer writer.Flush()

	for _, label := range client.Store.Labels {
		writer.Write([]string{IdFormat(label), "@" + label.Name})
	}

	return nil
}
