package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func List(config Config, sync lib.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	projectColorHash := GenerateColorHash(projectNames, colorList)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range sync.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DueDateTime(), item.AllDay),
			ProjectFormat(item, sync.Projects, projectColorHash),
			item.LabelsString(sync.Labels),
			item.Content,
		)
	}
	w.Flush()
	return nil
}
