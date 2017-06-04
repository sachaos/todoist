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
		return err
	}

	idCarrier, err := todoist.SearchByID(client.Store.Items, item_id)
	item := idCarrier.(todoist.Item)
	if err != nil {
		return err
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
		[]string{"Project", ProjectFormat(item.ProjectID, client.Store.Projects, projectColorHash, c)},
		[]string{"Labels", item.LabelsString(client.Store.Labels)},
		[]string{"Priority", PriorityFormat(item.Priority)},
		[]string{"DueDate", DueDateFormat(item.DueDateTime(), item.AllDay)},
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
