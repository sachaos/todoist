package main

import (
	"context"
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
	item.ProjectID = c.String("project-id")
	if item.ProjectID == "" {
		item.ProjectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}
	item.LabelNames = func(str string) []string {
		stringNames := strings.Split(str, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, stringName)
		}
		return names
	}(c.String("label-names"))

	item.DateString = c.String("date")
	item.AutoReminder = c.Bool("reminder")

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	return Sync(c)
}
