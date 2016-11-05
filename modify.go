package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

func Modify(config Config, sync lib.Sync, c *cli.Context) error {
	item := lib.Item{}
	next_project := lib.Project{}
	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item_id, err := strconv.Atoi(c.Args().First())
	item, err = sync.Items.FindByID(item_id)
	if err != nil {
		return err
	}
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

	next_project.ID = c.Int("project-id")

	if !c.Args().Present() {
		return CommandFailed
	}

	err = lib.UpdateItem(item, config.Token)
	if err != nil {
		return CommandFailed
	}

	err = lib.MoveItem(item, next_project, config.Token)
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(config, c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
