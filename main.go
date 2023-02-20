package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/rkoesters/xdg/basedir"
	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	homePath, _   = os.UserHomeDir()
	configPath    = filepath.Join(basedir.ConfigHome, "todoist")
	cachePath     = filepath.Join(basedir.CacheHome, "todoist", "cache.json")
	CommandFailed = errors.New("command failed")
	IdNotFound    = errors.New("specified id not found")
	writer        Writer
)

const (
	configName = "config"
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
	app.Version = "0.18.0"
	app.EnableBashCompletion = true

	contentFlag := cli.StringFlag{
		Name:  "content, c",
		Usage: "content",
	}
	priorityFlag := cli.IntFlag{
		Name:  "priority, p",
		Value: 4,
		Usage: "priority (1-4)",
	}
	labelNamesFlag := cli.StringFlag{
		Name:  "label-names, L",
		Usage: "label names (separated by ,)",
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
		Usage: "date string (today, 2020/04/02, 2020/03/21 18:00)",
	}
	browseFlag := cli.BoolFlag{
		Name:  "browse, o",
		Usage: "when contain URL, open it",
	}
	filterFlag := cli.StringFlag{
		Name:  "filter, f",
		Usage: "filter expression",
	}
	sortPriorityFlag := cli.BoolFlag{
		Name:  "priority, p",
		Usage: "sort the output by priority",
	}
	reminderFlg := cli.BoolFlag{
		Name:  "reminder, r",
		Usage: "set reminder (only premium users)",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "header",
			Usage: "output with header",
		},
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

		if err := LoadCache(cachePath, &store); err != nil {
			return err
		}

		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.AddConfigPath(".")
		viper.SetEnvPrefix("todoist") // uppercased automatically by viper
		viper.AutomaticEnv()

		var token string

		configFile := filepath.Join(configPath, configName+"."+configType)
		if err := AssureExists(configFile); err != nil {
			return err
		}

		if err := viper.ReadInConfig(); err != nil {
			if _, isConfigNotFoundError := err.(viper.ConfigFileNotFoundError); !isConfigNotFoundError {
				// config file was found but could not be read => not recoverable
				return err
			} else if !viper.IsSet("token") {
				// config file not found and token missing (not provided via another source,
				// such as environment variables) => ask interactively for token and store it in config file.
				fmt.Printf("Input API Token: ")
				fmt.Scan(&token)
				viper.Set("token", token)
				buf, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
				if err != nil {
					panic(fmt.Errorf("Fatal error config file: %s \n", err))
				}
				err = ioutil.WriteFile(configFile, buf, 0600)
				if err != nil {
					panic(fmt.Errorf("Fatal error config file: %s \n", err))
				}
			}
		}

		if exists, _ := Exists(configFile); exists {
			// Ensure that the config file has permission 0600, because it contains
			// the API token and should only be read by the user.
			// This is only necessary iff the config file exists, which may not be the case
			// when config is loaded from environment variables.
			fi, err := os.Lstat(configFile)
			if err != nil {
				panic(fmt.Errorf("Fatal error config file: %s \n", err))
			}
			if runtime.GOOS != "windows" && fi.Mode().Perm() != 0600 {
				panic(fmt.Errorf("Config file has wrong permissions. Make sure to give permissions 600 to file %s \n", configFile))
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
		} else {
			color.NoColor = false
		}

		if c.Bool("csv") {
			writer = csv.NewWriter(os.Stdout)
		} else if runtime.GOOS == "windows" && !color.NoColor {
			writer = NewTSVWriter(color.Output)
		} else {
			writer = NewTSVWriter(os.Stdout)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Show all tasks",
			Action:  List,
			Flags: []cli.Flag{
				filterFlag,
				sortPriorityFlag,
			},
			ArgsUsage: " ",
		},
		{
			Name:   "show",
			Usage:  "Show task detail",
			Action: Show,
			Flags: []cli.Flag{
				browseFlag,
			},
			ArgsUsage: "<Item ID>",
		},
		{
			Name:    "completed-list",
			Aliases: []string{"c-l", "cl"},
			Usage:   "Show all completed tasks (only premium user)",
			Action:  CompletedList,
			Flags: []cli.Flag{
				filterFlag,
			},
			ArgsUsage: " ",
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add task",
			Action:  Add,
			Flags: []cli.Flag{
				priorityFlag,
				labelNamesFlag,
				projectIDFlag,
				projectNameFlag,
				dateFlag,
				reminderFlg,
			},
			ArgsUsage: "<Item content>",
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify task",
			Action:  Modify,
			Flags: []cli.Flag{
				contentFlag,
				priorityFlag,
				labelNamesFlag,
				projectIDFlag,
				projectNameFlag,
				dateFlag,
			},
			ArgsUsage: "<Item ID>",
		},
		{
			Name:      "close",
			Aliases:   []string{"c"},
			Usage:     "Close task",
			Action:    Close,
			ArgsUsage: "<Item ID>",
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete task",
			Action:    Delete,
			ArgsUsage: "<Item ID>",
		},
		{
			Name:      "labels",
			Usage:     "Show all labels",
			Action:    Labels,
			ArgsUsage: " ",
		},
		{
			Name:      "projects",
			Usage:     "Show all projects",
			Action:    Projects,
			ArgsUsage: " ",
		},
		{
			Name:    "add-project",
			Aliases: []string{"ap"},
			Usage:   "Add new project",
			Action:  AddProject,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "color",
					Usage: "In range 30-49",
				},
				cli.IntFlag{
					Name:  "item-order",
					Usage: "Order index",
				},
			},
			ArgsUsage: "<Project name>",
		},
		{
			Name:      "karma",
			Usage:     "Show karma",
			Action:    Karma,
			ArgsUsage: " ",
		},
		{
			Name:      "sync",
			Aliases:   []string{"s"},
			Usage:     "Sync cache",
			Action:    Sync,
			ArgsUsage: " ",
		},
		{
			Name:      "quick",
			Aliases:   []string{"q"},
			Usage:     "Quick add a task",
			Action:    Quick,
			ArgsUsage: "<Item content>",
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
