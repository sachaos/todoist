package main

import (
	"errors"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

var (
	default_config_path = os.Getenv("HOME") + "/.todoist.config.json"
	default_cache_path  = os.Getenv("HOME") + "/.todoist.cache.json"
	CommandFailed       = errors.New("Command Failed")
)

const (
	ShortDateTimeFormat = "06/1/2(Mon) 15:04"
	ShortDateFormat     = "06/1/2(Mon)"
)

func main() {
	sync, err := LoadCache(default_cache_path)
	config, err := LoadConfig(default_config_path)
	if err != nil {
		return
	}
	app := cli.NewApp()
	app.Name = "todoist"
	app.Usage = "Todoist CLI Client"
	app.Version = "0.5.0"

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
	dateFlag := cli.StringFlag{
		Name:  "date, d",
		Usage: "date string (today, 2016/10/02, 2016/09/02 18:00)",
	}
	browseFlag := cli.BoolFlag{
		Name:  "browse, o",
		Usage: "when contain URL, open it",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "color",
		},
	}

	app.Before = func(c *cli.Context) error {
		if !c.Bool("color") {
			color.NoColor = true
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows all tasks",
			Action: func(c *cli.Context) error {
				return List(config, sync, c)
			},
		},
		{
			Name:  "show",
			Usage: "Show task detail",
			Action: func(c *cli.Context) error {
				return Show(config, sync, c)
			},
			Flags: []cli.Flag{
				browseFlag,
			},
		},
		{
			Name:  "completed-list",
			Usage: "Shows all completed tasks (only premium user)",
			Action: func(c *cli.Context) error {
				return CompletedList(config, sync, c)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add task",
			Action: func(c *cli.Context) error {
				return Add(config, sync, c)
			},
			Flags: []cli.Flag{
				priorityFlag,
				labelIDsFlag,
				projectIDFlag,
				dateFlag,
			},
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify task",
			Action: func(c *cli.Context) error {
				return Modify(config, sync, c)
			},
			Flags: []cli.Flag{
				contentFlag,
				priorityFlag,
				labelIDsFlag,
				projectIDFlag,
				dateFlag,
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
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete task",
			Action: func(c *cli.Context) error {
				return Delete(config, c)
			},
		},
		{
			Name:  "labels",
			Usage: "Shows all labels",
			Action: func(c *cli.Context) error {
				return Labels(config, sync, c)
			},
		},
		{
			Name:  "projects",
			Usage: "Shows all projects",
			Action: func(c *cli.Context) error {
				return Projects(config, sync, c)
			},
		},
		{
			Name:  "karma",
			Usage: "Show karma",
			Action: func(c *cli.Context) error {
				return Karma(config, sync, c)
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
