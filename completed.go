package main

import (
	"context"

	"github.com/sachaos/todoist/lib"

	"github.com/urfave/cli"
)

func CompletedList(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []int
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	var completed todoist.Completed

	if err := client.CompletedAll(context.Background(), &completed); err != nil {
		return err
	}

	defer writer.Flush()

	for _, item := range completed.Items {
		writer.Write([]string{
			IdFormat(item),
			CompletedDateFormat(item.CompletedDateTime()),
			ProjectFormat(item.ProjectID, client.Store.Projects, projectColorHash, c),
			ContentFormat(item),
		})
	}

	return nil
}
