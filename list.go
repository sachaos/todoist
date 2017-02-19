package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func List(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	itemList := [][]string{}
	for _, itemOrder := range sync.ItemOrders {
		item := itemOrder.Data.(todoist.Item)
		if item.Checked == 1 {
			continue
		}
		itemList = append(itemList, []string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DueDateTime(), item.AllDay),
			ProjectFormat(item.ProjectID, sync.Projects, projectColorHash, c),
			item.LabelsString(sync.Labels),
			ContentPrefix(sync.Items, item, c) + ContentFormat(item),
		})
	}

	defer writer.Flush()

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
