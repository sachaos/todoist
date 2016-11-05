package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Sync(config Config, c *cli.Context) (lib.Sync, error) {
	var sync lib.Sync
	sync, err := lib.SyncAll(config.Token)
	if err != nil {
		return sync, CommandFailed
	}
	err = WriteCache(default_cache_path, sync)
	if err != nil {
		return sync, CommandFailed
	}

	return sync, nil
}
