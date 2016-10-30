package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

func Modify(config Config, c *cli.Context) error {
	item := lib.Item{}
	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item.ID, err = strconv.Atoi(c.Args().First())
	item.Content = c.String("content")
	item.Priority = c.Int("priority")
	item.ProjectID = c.Int("project-id")
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

	if !c.Args().Present() {
		return CommandFailed
	}

	err = lib.UpdateItem(item, config.Token)
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(config, c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
