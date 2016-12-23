package main

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"text/tabwriter"
)

func Show(config Config, sync lib.Sync, c *cli.Context) error {
	item_id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return err
	}

	item, err := sync.Items.FindByID(item_id)
	if err != nil {
		return err
	}

	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	projectColorHash := GenerateColorHash(projectNames, colorList)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	fmt.Fprintf(w, "ID\t%s\n", IdFormat(item))
	fmt.Fprintf(w, "Content\t%s\n", ContentFormat(item))
	fmt.Fprintf(w, "Project\t%s\n", ProjectFormat(item, sync.Projects, projectColorHash))
	fmt.Fprintf(w, "Labels\t%s\n", item.LabelsString(sync.Labels))
	fmt.Fprintf(w, "Priority\t%s\n", PriorityFormat(item.Priority))
	fmt.Fprintf(w, "DueDate\t%s\n", DueDateFormat(item.DueDateTime(), item.AllDay))
	if lib.HasURL(item) {
		fmt.Fprintf(w, "URL\t%s\n", lib.GetContentURL(item))
		if c.Bool("browse") {
			browser.OpenURL(lib.GetContentURL(item))
		}
	}

	w.Flush()
	return nil
}
