package main

import (
	"context"
	"fmt"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func AddSection(c *cli.Context) error {
	client := GetClient(c)

	section := todoist.Section{}
	if !c.Args().Present() {
		return CommandFailed
	}

	section.Name = c.Args().First()

	projectName := c.String("project-name")
	if projectName != "" {
		projectID := client.Store.Projects.GetIDByName(projectName)
		if projectID == "" {
			return fmt.Errorf("Did not find a project named '%v'", projectName)
		}
		section.ProjectID = projectID
	} else {
		section.ProjectID = c.String("project-id")
	}

	if err := client.AddSection(context.Background(), section); err != nil {
		return err
	}

	return Sync(c)
}
