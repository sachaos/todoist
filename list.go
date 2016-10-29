package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func List(config Config, c *cli.Context) error {
	var sync lib.Sync

	sync, err := lib.LoadCache(default_cache_path)

	if err != nil {
		sync, err = Sync(config, c)
		if err != nil {
			return err
		}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range sync.Items {
		fmt.Fprintf(w, "%d\tp%d\t%s\t%s\t%s\n",
			item.ID,
			item.Priority,
			item.ProjectString(sync.Projects),
			item.LabelsString(sync.Labels),
			item.Content,
		)
	}
	w.Flush()
	return nil
}
