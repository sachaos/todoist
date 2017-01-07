package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func Projects(sync todoist.Sync, c *cli.Context) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, project := range sync.Projects {
		fmt.Fprintf(w, "%d\t%s\n", project.ID, "#"+project.Name)
	}
	w.Flush()
	return nil
}
