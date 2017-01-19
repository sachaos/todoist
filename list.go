package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func List(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	projectColorHash := GenerateColorHash(projectNames, colorList)

	defer writer.Flush()

	for _, item := range sync.Items {
		writer.Write([]string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DueDateTime(), item.AllDay),
			ProjectFormat(item, sync.Projects, projectColorHash),
			item.LabelsString(sync.Labels),
			ContentFormat(item),
		})
	}

	return nil
}
