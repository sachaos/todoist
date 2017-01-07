package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func CompletedList(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	completed, err := todoist.CompletedAll(viper.GetString("token"))
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
