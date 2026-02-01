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

	item.LabelNames = func(str string) []string {
		stringNames := strings.Split(str, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, stringName)
		}
		return names
	}(c.String("label-names"))

	item.Due = &todoist.Due{String: c.String("date")}

	// Note: AutoReminder via Sync API doesn't work reliably with API v1
	// We'll create a proper reminder via REST API v2 after task creation
	wantReminder := c.Bool("reminder")

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	// Sync to get the created task's ID
	if err := Sync(c); err != nil {
		return err
	}

	// If reminder was requested and we have a due date, create a reminder
	if wantReminder && item.Due != nil && item.Due.String != "" {
		// Find the newly created task by matching content
		var createdItemID string
		for _, storedItem := range client.Store.Items {
			if storedItem.Content == item.Content {
				createdItemID = storedItem.ID
				break
			}
		}

		if createdItemID != "" {
			// Create a reminder at the due date using Sync API
			err := client.AddReminder(context.Background(), createdItemID, item.Due)
			if err != nil {
				// Log warning but don't fail the whole operation
				fmt.Printf("Warning: task created but reminder failed: %v\n", err)
			}
		}
	}

	return nil
}
