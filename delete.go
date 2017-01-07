package main

import (
	"strconv"

	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func Delete(c *cli.Context) error {
	item_ids := []int{}
	for _, arg := range c.Args() {
		item_id, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		item_ids = append(item_ids, item_id)
	}

	if len(item_ids) == 0 {
		return CommandFailed
	}

	err := todoist.DeleteItem(item_ids, viper.GetString("token"))
	if err != nil {
		return CommandFailed
	}

	_, err = Sync(c)
	if err != nil {
		return CommandFailed
	}

	return nil
}
