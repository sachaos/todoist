package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
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
	if c.Args().Len() != 1 {
		return fmt.Errorf("add command requires 1 positional argument for the task title, but got %v.", c.Args().Len())
	}

	item.Content = c.Args().First()

	item.Priority = priorityMapping[c.Int("priority")]

	projectName := c.String("project-name")
	if projectName != "" {
		projectId := client.Store.Projects.GetIDByName(projectName)
		if projectId == "" {
			return fmt.Errorf("Did not find a project named '%v'", projectName)
		}
		item.ProjectID = projectId
	} else {
		item.ProjectID = c.String("project-id")
	}

	labelNames := c.String("label-names")
	if labelNames != "" {
		stringNames := strings.Split(labelNames, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, strings.TrimSpace(stringName))
		}
		item.LabelNames = names
	}

	sectionName := c.String("section-name")
	if sectionName != "" {
		sectionID := client.Store.Sections.GetIDByName(sectionName, item.ProjectID)
		if sectionID == "" {
			return fmt.Errorf("Did not find a section named '%v'", sectionName)
		}
		item.SectionID = sectionID
	} else if c.String("section-id") != "" {
		item.SectionID = c.String("section-id")
	}

	item.Description = c.String("description")

	item.Due = &todoist.Due{String: c.String("date")}

	if c.IsSet("deadline") {
		deadlineDate := c.String("deadline")
		if deadlineDate == "" || deadlineDate == "null" {
			item.Deadline = &todoist.Deadline{Date: ""}
		} else {
			item.Deadline = &todoist.Deadline{Date: deadlineDate}
		}
	}

	item.AutoReminder = c.Bool("reminder")

	newID, err := client.AddItem(context.Background(), item)
	if err != nil {
		return err
	}

	if c.IsSet("deadline") {
		deadlineDate := c.String("deadline")
		if err := client.UpdateItemDeadline(context.Background(), newID, deadlineDate); err != nil {
			return fmt.Errorf("failed to update deadline for task %s: %w", newID, err)
		}
	}

	return Sync(c)
}
