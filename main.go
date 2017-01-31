package main

import (
	"errors"
	"fmt"
	"os"

	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	configPath         = os.Getenv("HOME")
	default_cache_path = os.Getenv("HOME") + "/.todoist.cache.json"
	CommandFailed      = errors.New("Command Failed")
	writer             Writer
)

const (
	configName = ".todoist.config"
	configType = "json"

	ShortDateTimeFormat = "06/1/2(Mon) 15:04"
	ShortDateFormat     = "06/1/2(Mon)"
)

func main() {
	sync, err := LoadCache(default_cache_path)
	if err != nil {
		return
	}

	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()

	if err != nil {
		var token string
		fmt.Printf("Input API Token: ")
		fmt.Scan(&token)
		viper.Set("token", token)
		buf, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		err = ioutil.WriteFile(filepath.Join(configPath, configName+"."+configType), buf, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	app := cli.NewApp()
	app.Name = "todoist"
	app.Usage = "Todoist CLI Client"
	app.Version = "0.6.0"

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
		cli.BoolFlag{
			Name: "csv",
		},
		cli.BoolFlag{
			Name: "namespace",
		},
		cli.BoolFlag{
			Name: "indent",
		},
		cli.BoolFlag{
			Name: "project-namespace",
		},
	}

	app.Before = func(c *cli.Context) error {
		if !c.Bool("color") {
			color.NoColor = true
		}

		if c.GlobalBool("csv") {
			writer = csv.NewWriter(os.Stdout)
		} else {
			writer = NewTSVWriter(os.Stdout)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows all tasks",
			Action: func(c *cli.Context) error {
				return List(sync, c)
			},
		},
		{
			Name:  "show",
			Usage: "Show task detail",
			Action: func(c *cli.Context) error {
				return Show(sync, c)
			},
			Flags: []cli.Flag{
				browseFlag,
			},
		},
		{
			Name:  "completed-list",
			Usage: "Shows all completed tasks (only premium user)",
			Action: func(c *cli.Context) error {
				return CompletedList(sync, c)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add task",
			Action: func(c *cli.Context) error {
				return Add(sync, c)
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
				return Modify(sync, c)
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
				return Close(c)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete task",
			Action: func(c *cli.Context) error {
				return Delete(c)
			},
		},
		{
			Name:  "labels",
			Usage: "Shows all labels",
			Action: func(c *cli.Context) error {
				return Labels(sync, c)
			},
		},
		{
			Name:  "projects",
			Usage: "Shows all projects",
			Action: func(c *cli.Context) error {
				return Projects(sync, c)
			},
		},
		{
			Name:  "karma",
			Usage: "Show karma",
			Action: func(c *cli.Context) error {
				return Karma(sync, c)
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync cache",
			Action: func(c *cli.Context) error {
				_, err := Sync(c)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
