package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func Projects(sync todoist.Sync, c *cli.Context) error {
	defer writer.Flush()

	for _, project := range sync.Projects {
		writer.Write([]string{IdFormat(project), "#" + project.Name})
	}

	return nil
}
