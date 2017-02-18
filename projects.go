package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Projects(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	itemList := [][]string{}
	for _, itemOrder := range sync.ProjectOrders {
		project := itemOrder.Data.(todoist.Project)
		itemList = append(itemList, []string{IdFormat(project), ProjectFormat(project.ID, sync.Projects, projectColorHash, c)})
	}

	defer writer.Flush()

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
