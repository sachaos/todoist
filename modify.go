package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
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
	item.Description = c.String("description")
	item.Priority = priorityMapping[c.Int("priority")]
	if labelNames := c.String("label-names"); labelNames != "" {
		stringNames := strings.Split(labelNames, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, strings.TrimSpace(stringName))
		}
		item.LabelNames = names
	}

	item.Due = &todoist.Due{String: c.String("date")}

	projectID := c.String("project-id")
	if projectID == "" {
		projectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}

	sectionName := c.String("section-name")
	sectionID := c.String("section-id")
	if sectionName != "" {
		sectionID = client.Store.Sections.GetIDByName(sectionName, projectID)
		if sectionID == "" {
			return fmt.Errorf("Did not find a section named '%v'", sectionName)
		}
	}

	if err := client.UpdateItem(context.Background(), *item); err != nil {
		return err
	}

	if projectID != "" {
		if err := client.MoveItem(context.Background(), item, projectID); err != nil {
			return err
		}
	}

	if sectionID != "" {
		if err := client.MoveItemToSection(context.Background(), item, sectionID); err != nil {
			return err
		}
	}

	if c.IsSet("deadline") {
		deadlineDate := c.String("deadline")
		if err := client.UpdateItemDeadline(context.Background(), item_id, deadlineDate); err != nil {
			return fmt.Errorf("failed to update deadline for task %s: %w", item_id, err)
		}
	}

	return Sync(c)
}
