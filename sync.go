package main

import (
	"context"
	"fmt"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Sync(c *cli.Context) error {
	client := GetClient(c)

	pc, err := GetPipelineCache(pipelineCachePath)
	if err == nil && !pc.IsEmpty() {
		fmt.Println("Syncing pending tasks...")

		items := pc.GetItems()
		commands := todoist.Commands{}
		uuids := []string{}

		for _, pipelineItem := range items {
			if pipelineItem.IsQuick {
				err := client.QuickCommand(context.Background(), pipelineItem.QuickText)
				if err != nil {
					fmt.Printf("Warning: failed to sync quick command: %v\n", err)
					continue
				}
				uuids = append(uuids, pipelineItem.Command.UUID)
			} else {
				commands = append(commands, pipelineItem.Command)
				uuids = append(uuids, pipelineItem.Command.UUID)
			}
		}

		if len(commands) > 0 {
			err = client.ExecCommands(context.Background(), commands)
			if err != nil {
				return fmt.Errorf("failed to sync pending tasks: %w", err)
			}
		}

		if len(uuids) > 0 {
			err = pc.RemoveItems(uuids)
			if err != nil {
				fmt.Printf("Warning: failed to remove items from pipeline cache: %v\n", err)
			} else {
				err = WritePipelineCache(pipelineCachePath, pc)
				if err != nil {
					fmt.Printf("Warning: failed to write pipeline cache: %v\n", err)
				}
			}
		}
	}

	err = client.Sync(context.Background())
	if err != nil {
		return err
	}
	return WriteCache(cachePath, client.Store)
}
