package main

import (
	"strconv"
	"strings"

	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func Add(sync lib.Sync, c *cli.Context) error {
	item := lib.Item{}
	if !c.Args().Present() {
		return CommandFailed
	}

	item.Content = c.Args().First()
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

	item.DateString = c.String("date")

	err := lib.AddItem(item, viper.GetString("token"))
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
