package main

import (
	"context"
	todoist "github.com/sachaos/todoist/lib"
	"strings"

	"github.com/urfave/cli"
)

func Modify(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item_id, err := client.CompleteItemIDByPrefix(c.Args().First())
	if err != nil {
		return err
	}
	item := client.Store.FindItem(item_id)
	if item == nil {
		return IdNotFound
	}
	item.Content = c.String("content")
	item.Priority = priorityMapping[c.Int("priority")]
	item.LabelNames = func(str string) []string {
		stringNames := strings.Split(str, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, stringName)
		}
		return names
	}(c.String("label-names"))

	if date := c.String("date"); date != "" {
		item.Due = &todoist.Due{Date: date}
	}

	projectID := c.String("project-id")
	if projectID == "" {
		projectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}

	if !c.Args().Present() {
		return CommandFailed
	}

	if err := client.UpdateItem(context.Background(), *item); err != nil {
		return err
	}

	if err := client.MoveItem(context.Background(), item, projectID); err != nil {
		return err
	}

	return Sync(c)
}
