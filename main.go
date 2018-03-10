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
	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	configPath         = os.Getenv("HOME")
	default_cache_path = os.Getenv("HOME") + "/.todoist.cache.json"
	CommandFailed      = errors.New("command failed")
	IdNotFound         = errors.New("specified id not found")
	writer             Writer
)

const (
	configName = ".todoist.config"
	configType = "json"

	ShortDateTimeFormat = "06/01/02(Mon) 15:04"
	ShortDateFormat     = "06/01/02(Mon)"
)

func GetClient(c *cli.Context) *todoist.Client {
	return c.App.Metadata["client"].(*todoist.Client)
}

func main() {
	app := cli.NewApp()
	app.Name = "todoist"
	app.Usage = "Todoist CLI Client"
	app.Version = "0.10.1"

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
	projectNameFlag := cli.StringFlag{
		Name:  "project-name, N",
		Usage: "project name",
	}
	dateFlag := cli.StringFlag{
		Name:  "date, d",
		Usage: "date string (today, 2016/10/02, 2016/09/02 18:00)",
	}
	browseFlag := cli.BoolFlag{
		Name:  "browse, o",
		Usage: "when contain URL, open it",
	}
	filterFlag := cli.StringFlag{
		Name:  "filter, f",
		Usage: "filter expression",
	}
	reminderFlg := cli.BoolFlag{
		Name:  "reminder, r",
		Usage: "set reminder (only premium user)",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "color",
			Usage: "colorize output",
		},
		cli.BoolFlag{
			Name:  "csv",
			Usage: "output in CSV format",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "output logs",
		},
		cli.BoolFlag{
			Name:  "namespace",
			Usage: "display parent task like namespace",
		},
		cli.BoolFlag{
			Name:  "indent",
			Usage: "display children task with indent",
		},
		cli.BoolFlag{
			Name:  "project-namespace",
			Usage: "display parent project like namespace",
		},
	}

	app.Before = func(c *cli.Context) error {
		var store todoist.Store

		if err := LoadCache(default_cache_path, &store); err != nil {
			return err
		}

		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.AddConfigPath(".")

		var token string

		if err := viper.ReadInConfig(); err != nil {
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

		config := &todoist.Config{AccessToken: viper.GetString("token"), DebugMode: c.Bool("debug"), Color: viper.GetBool("color")}

		client := todoist.NewClient(config)
		client.Store = &store

		app.Metadata = map[string]interface{}{
			"client": client,
			"config": config,
		}

		if !c.Bool("color") && !config.Color {
			color.NoColor = true
		}

		if c.Bool("csv") {
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
			Action:  List,
			Flags: []cli.Flag{
				filterFlag,
			},
		},
		{
			Name:   "show",
			Usage:  "Show task detail",
			Action: Show,
			Flags: []cli.Flag{
				browseFlag,
			},
		},
		{
			Name:    "completed-list",
			Aliases: []string{"c-l", "cl"},
			Usage:   "Shows all completed tasks (only premium user)",
			Action:  CompletedList,
			Flags: []cli.Flag{
				filterFlag,
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add task",
			Action:  Add,
			Flags: []cli.Flag{
				priorityFlag,
				labelIDsFlag,
				projectIDFlag,
				projectNameFlag,
				dateFlag,
				reminderFlg,
			},
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify task",
			Action:  Modify,
			Flags: []cli.Flag{
				contentFlag,
				priorityFlag,
				labelIDsFlag,
				projectIDFlag,
				projectNameFlag,
				dateFlag,
			},
		},
		{
			Name:    "close",
			Aliases: []string{"c"},
			Usage:   "Close task",
			Action:  Close,
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete task",
			Action:  Delete,
		},
		{
			Name:   "labels",
			Usage:  "Shows all labels",
			Action: Labels,
		},
		{
			Name:   "projects",
			Usage:  "Shows all projects",
			Action: Projects,
		},
		{
			Name:   "karma",
			Usage:  "Show karma",
			Action: Karma,
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync cache",
			Action:  Sync,
		},
		{
			Name:    "quick",
			Aliases: []string{"q"},
			Usage:   "Quick add a task",
			Action:  Quick,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
