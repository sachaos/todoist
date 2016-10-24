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
	config              = Config{}
)

func main() {
	err := Setup(default_config_path)
	if err != nil {
		return
	}
	app := cli.NewApp()
	app.Name = "todoist"
	app.Usage = "Todoist cli client"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows all tasks",
			Action:  List,
		},
		// {
		// 	Name:    "add",
		// 	Aliases: []string{"a"},
		// 	Usage:   "Add task",
		// 	Action:  Add,
		// },
	}
	app.Run(os.Args)
}
