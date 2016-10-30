package main

import (
	"errors"
	"github.com/urfave/cli"
	"os"
)

var (
	default_config_path = os.Getenv("HOME") + "/.todoist.config.json"
	default_cache_path  = os.Getenv("HOME") + "/.todoist.cache.json"
	CommandFailed       = errors.New("Command Failed")
)

func main() {
	config, err := Setup(default_config_path)
	if err != nil {
		return
	}
	app := cli.NewApp()
	app.Name = "todoist"
	app.Usage = "Todoist CLI Client"
	app.Version = "0.1.0"

	contentFlag := cli.StringFlag{
		Name:  "content, c",
		Usage: "content",
	}
	priorityFlag := cli.IntFlag{
		Name:  "priority, p",
		Value: 1,
		Usage: "priority (1-4)",
	}
	labelIDsFlag := cli.StringFlag{
		Name:  "label-ids, L",
		Usage: "label ids (separated by ,)",
	}
	projectIDFlag := cli.IntFlag{
		Name:  "project-id, P",
		Usage: "project id",
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows all tasks",
			Action: func(c *cli.Context) error {
				return List(config, c)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add task",
			Action: func(c *cli.Context) error {
				return Add(config, c)
			},
			Flags: []cli.Flag{
				priorityFlag,
				labelIDsFlag,
				projectIDFlag,
			},
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify task",
			Action: func(c *cli.Context) error {
				return Modify(config, c)
			},
			Flags: []cli.Flag{
				contentFlag,
				priorityFlag,
				labelIDsFlag,
			},
		},
		{
			Name:    "close",
			Aliases: []string{"c"},
			Usage:   "Close task",
			Action: func(c *cli.Context) error {
				return Close(config, c)
			},
		},
		{
			Name:  "labels",
			Usage: "Shows all labels",
			Action: func(c *cli.Context) error {
				return Labels(config, c)
			},
		},
		{
			Name:  "projects",
			Usage: "Shows all projects",
			Action: func(c *cli.Context) error {
				return Projects(config, c)
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync cache",
			Action: func(c *cli.Context) error {
				_, err := Sync(config, c)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
