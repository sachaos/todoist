package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func CompletedList(config Config, sync lib.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	completed, err := lib.CompletedAll(config.Token)
	if err != nil {
		return err
	}

	projectColorHash := GenerateColorHash(projectNames, colorList)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, item := range completed.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			IdFormat(item),
			CompletedDateFormat(item.CompletedDateTime()),
			ProjectFormat(item, sync.Projects, projectColorHash),
			ContentFormat(item),
		)
	}
	w.Flush()
	return nil
}
