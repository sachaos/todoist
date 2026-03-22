package main

import (
	"strings"

	"github.com/pkg/browser"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Show(c *cli.Context) error {
	client := GetClient(c)

	id := c.Args().First()

	// Check if the ID refers to a section
	section := client.Store.FindSection(id)
	if section != nil {
		return showSection(c, client, section)
	}

	item := client.Store.FindItem(id)
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
		[]string{"ID", IdFormat(item)},
		[]string{"Content", ContentFormat(item)},
		[]string{"Project", ProjectFormat(item.ProjectID, client.Store, projectColorHash, c)},
		[]string{"Section", SectionFormat(item.SectionID, client.Store, c)},
		[]string{"Labels", item.LabelsString()},
		[]string{"Priority", PriorityFormat(item.Priority)},
		[]string{"DueDate", DueDateFormat(item.DateTime(), item.AllDay)},
		[]string{"URL", strings.Join(todoist.GetContentURL(item), ",")},
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

func showSection(c *cli.Context, client *todoist.Client, section *todoist.Section) error {
	colorList := ColorList()
	var projectIds []string
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	records := [][]string{
		{"ID", IdFormat(section)},
		{"Name", section.Name},
		{"Project", ProjectFormat(section.ProjectID, client.Store, projectColorHash, c)},
	}
	defer writer.Flush()

	for _, record := range records {
		writer.Write(record)
	}
	return nil
}
