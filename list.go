package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
	"time"
)

func DueDateString(item lib.Item) string {
	due_date := item.DueDateTime()
	if (due_date == time.Time{}) {
		return ""
	}
	due_date = due_date.Local()
	if !item.AllDay {
		return due_date.Format(ShortDateTimeFormat)
	}
	return due_date.Format(ShortDateFormat)
}

func List(config Config, sync lib.Sync, c *cli.Context) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range sync.Items {
		fmt.Fprintf(w, "%d\tp%d\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.Priority,
			DueDateString(item),
			item.ProjectString(sync.Projects),
			item.LabelsString(sync.Labels),
			item.Content,
		)
	}
	w.Flush()
	return nil
}
