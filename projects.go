package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Projects(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []int
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	itemList := [][]string{}
	for _, itemOrder := range client.Store.ProjectOrders {
		project := itemOrder.Data.(todoist.Project)
		itemList = append(itemList, []string{IdFormat(project), ProjectFormat(project.ID, client.Store.Projects, projectColorHash, c)})
	}

	defer writer.Flush()

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
