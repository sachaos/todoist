package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

var priorityMapping = map[int]int{
	1: 4,
	2: 3,
	3: 2,
	4: 1,
}

func Add(c *cli.Context) error {
	client := GetClient(c)

	item := todoist.Item{}
	if !c.Args().Present() {
		return CommandFailed
	}

	item.Content = c.Args().First()
	item.Priority = priorityMapping[c.Int("priority")]
	item.ProjectID = c.Int("project-id")
	if item.ProjectID == 0 {
		item.ProjectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}
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
	item.AutoReminder = c.Bool("reminder")

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	return Sync(c)
}
