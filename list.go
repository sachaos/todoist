package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []int
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	itemList := [][]string{}
	for _, itemOrder := range client.Store.ItemOrders {
		item := itemOrder.Data.(todoist.Item)
		r, err := Eval(ex, item)
		if err != nil {
			return err
		}
		if !r || item.Checked == 1 {
			continue
		}
		itemList = append(itemList, []string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DueDateTime(), item.AllDay),
			ProjectFormat(item.ProjectID, client.Store.Projects, projectColorHash, c),
			item.LabelsString(client.Store.Labels),
			ContentPrefix(client.Store.Items, item, c) + ContentFormat(item),
		})
	}

	defer writer.Flush()

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
