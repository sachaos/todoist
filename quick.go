package main

import (
	"fmt"
	"time"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Quick(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	quickText := c.Args().First()

	pc, err := GetPipelineCache(pipelineCachePath)
	if err != nil {
		return fmt.Errorf("failed to get pipeline cache: %w", err)
	}

	command := todoist.NewCommand("quick_add", map[string]interface{}{
		"text": quickText,
	})

	pipelineItem := PipelineItem{
		Item:      todoist.Item{},
		Command:   command,
		QuickText: quickText,
		IsQuick:   true,
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
