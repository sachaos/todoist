package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Modify(c *cli.Context) error {
	client := GetClient(c)

	next_project := todoist.Project{}
	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item_id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return err
	}
	idCarrier, err := todoist.SearchByID(client.Store.Items, item_id)
	if err != nil {
		return IdNotFound
	}
	item := idCarrier.(todoist.Item)
	item.Content = c.String("content")
	item.Priority = c.Int("priority")
	item.LabelIDs = func(str string) []int {
		stringIDs := strings.Split(str, ",")
		ids := []int{}
		for _, stringID := range stringIDs {
			id, err := strconv.Atoi(stringID)
			if err != nil {
				continue
			}
			ids = append(ids, id)
		}
		return ids
	}(c.String("label-ids"))

	item.DateString = c.String("date")

	projectID := c.Int("project-id")
	if projectID == 0 {
		projectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}
	next_project.ID = projectID

	if !c.Args().Present() {
		return CommandFailed
	}

	if err := client.UpdateItem(context.Background(), item); err != nil {
		return err
	}

	if err := client.MoveItem(context.Background(), item, next_project); err != nil {
		return err
	}

	if err := Sync(c); err != nil {
		return err
	}

	return nil
}
