package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func CompletedList(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectNames := []string{}
	for _, project := range sync.Projects {
		projectNames = append(projectNames, project.Name)
	}
	completed, err := todoist.CompletedAll(viper.GetString("token"))
	if err != nil {
		return err
	}

	projectColorHash := GenerateColorHash(projectNames, colorList)

	defer writer.Flush()

	for _, item := range completed.Items {
		writer.Write([]string{
			IdFormat(item),
			CompletedDateFormat(item.CompletedDateTime()),
			ProjectFormat(item, sync.Projects, projectColorHash),
			ContentFormat(item),
		})
	}

	return nil
}
