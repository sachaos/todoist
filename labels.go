package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func Labels(config Config, sync lib.Sync, c *cli.Context) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, label := range sync.Labels {
		fmt.Fprintf(w, "%d\t%s\n", label.ID, "@"+label.Name)
	}
	w.Flush()
	return nil
}
