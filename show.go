package main

import (
	"strings"

	"github.com/pkg/browser"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Show(c *cli.Context) error {
	client := GetClient(c)

	item_id := c.Args().First()

	item := client.Store.FindItem(item_id)
	if item == nil {
		return IdNotFound
	}

	colorList := ColorList()
	var projectIds []string
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	records := [][]string{
		{"ID", IdFormat(item)},
		{"Content", ContentFormat(item)},
		{"Description", DescriptionFormat(item)},
		{"Project", ProjectFormat(item.ProjectID, client.Store, projectColorHash, c)},
		{"Labels", item.LabelsString(client.Store)},
		{"Priority", PriorityFormat(item.Priority)},
		{"DueDate", DueDateFormat(item.DateTime(), item.AllDay)},
		{"URL", strings.Join(todoist.GetContentURL(item), ",")},
	}
	defer writer.Flush()

	for _, record := range records {
		writer.Write(record)
	}

	if todoist.HasURL(item) {
		if c.Bool("browse") {
			for _, url := range todoist.GetContentURL(item) {
				browser.OpenURL(url)
			}
		}
	}
	return nil
}
