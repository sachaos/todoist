package main

import (
	"fmt"
	"strings"
	"time"

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

	item.AutoReminder = c.Bool("reminder")

	pc, err := GetPipelineCache(pipelineCachePath)
	if err != nil {
		return fmt.Errorf("failed to get pipeline cache: %w", err)
	}

	command := todoist.NewCommand("item_add", item.AddParam())

	pipelineItem := PipelineItem{
		Item:      item,
		Command:   command,
		IsQuick:   false,
		CreatedAt: time.Now(),
	}

	err = pc.AddItem(pipelineItem)
	if err != nil {
		return fmt.Errorf("failed to add item to pipeline cache: %w", err)
	}

	err = WritePipelineCache(pipelineCachePath, pc)
	if err != nil {
		return fmt.Errorf("failed to write pipeline cache: %w", err)
	}

	fmt.Println("Task added to queue for syncing")
	StartBackgroundSync(client, pipelineCachePath, cachePath)
	return nil
}
