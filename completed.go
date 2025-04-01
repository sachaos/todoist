package main

import (
	"context"

	todoist "github.com/sachaos/todoist/lib"

	"github.com/urfave/cli/v2"
)

func CompletedList(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []string
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	var completed todoist.Completed

	if err := client.CompletedAll(context.Background(), &completed); err != nil {
		return err
	}

	defer writer.Flush()

	if c.Bool("header") {
		writer.Write([]string{"ID", "CompletedDate", "Project", "Content"})
	}

	for _, item := range completed.Items {
		result, err := Eval(ex, item, client.Store.Projects, client.Store.Sections, client.Store.Labels)
		if err != nil {
			return err
		}
		if !result {
			continue
		}
		writer.Write([]string{
			IdFormat(item),
			CompletedDateFormat(item.DateTime()),
			ProjectFormat(item.ProjectID, client.Store, projectColorHash, c),
			ContentFormat(item),
		})
	}

	return nil
}
