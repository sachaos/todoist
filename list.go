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
		sync, err = lib.FetchCache(config.Token)
		if err != nil {
			return CommandFailed
		}
		err = lib.SaveCache(default_cache_path, sync)
		if err != nil {
			return CommandFailed
		}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range sync.Items {
		fmt.Fprintf(w, "%d\tp%d\t%s\t%s\n", item.ID, item.Priority, lib.LabelsString(item, sync.Labels), item.Content)
		// for _, label_id := range item.LabelIDs {
		// 	label, err := FindByID(sync.Labels, label_id)
		// 	if err != nil {
		// 		return
		// 	}
		// 	fmt.Printf("@%s", label.Name)
		// }
		// fmt.Printf("\n")
	}
	w.Flush()
	return nil
}
