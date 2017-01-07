package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Karma(sync todoist.Sync, c *cli.Context) error {
	fmt.Println(sync.User.Karma)
	return nil
}
