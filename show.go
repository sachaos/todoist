package main

import (
	"strconv"

	"github.com/pkg/browser"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Show(c *cli.Context) error {
	client := GetClient(c)

	item_id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return CommandFailed
	}

	item := client.Store.FindItem(item_id)
	if item == nil {
		return IdNotFound
	}

	colorList := ColorList()
	var projectIds []int
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	records := [][]string{
		[]string{"ID", IdFormat(item)},
		[]string{"Content", ContentFormat(item)},
		[]string{"Project", ProjectFormat(item.ProjectID, client.Store, projectColorHash, c)},
		[]string{"Labels", item.LabelsString(client.Store)},
		[]string{"Priority", PriorityFormat(item.Priority)},
		[]string{"DueDate", DueDateFormat(item.DateTime(), item.AllDay)},
		[]string{"URL", todoist.GetContentURL(item)},
	}
	defer writer.Flush()

	for _, record := range records {
		writer.Write(record)
	}

	if todoist.HasURL(item) {
		if c.Bool("browse") {
			browser.OpenURL(todoist.GetContentURL(item))
		}
	}
	return nil
}
