package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func Sync(c *cli.Context) (todoist.Sync, error) {
	var sync todoist.Sync
	sync, err := todoist.SyncAll(viper.GetString("token"))
	if err != nil {
		return sync, err
	}
	err = WriteCache(default_cache_path, sync)
	if err != nil {
		return sync, err
	}

	return sync, nil
}
