package main

import (
	"strconv"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func Modify(sync todoist.Sync, c *cli.Context) error {
	next_project := todoist.Project{}
	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item_id, err := strconv.Atoi(c.Args().First())
	idCarrier, err := todoist.SearchByID(sync.Items, item_id)
	item := idCarrier.(todoist.Item)
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

	item.DateString = c.String("date")

	next_project.ID = c.Int("project-id")

	if !c.Args().Present() {
		return CommandFailed
	}

	err = todoist.UpdateItem(item, viper.GetString("token"))
	if err != nil {
		return CommandFailed
	}

	err = todoist.MoveItem(item, next_project, viper.GetString("token"))
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
