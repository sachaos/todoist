package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func CompletedList(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectTree := ProjectTree(sync)
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	completed, err := todoist.CompletedAll(viper.GetString("token"))
	if err != nil {
		return err
	}

	defer writer.Flush()

	for _, item := range completed.Items {
		writer.Write([]string{
			IdFormat(item),
			CompletedDateFormat(item.CompletedDateTime()),
			ProjectFormat(item.ProjectID, projectTree, projectColorHash, c),
			ContentFormat(item),
		})
	}

	return nil
}
