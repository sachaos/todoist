package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func MoveSection(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	sectionID := c.Args().First()
	if client.Store.FindSection(sectionID) == nil {
		return IdNotFound
	}

	projectName := c.String("project-name")
	projectID := c.String("project-id")
	if projectName != "" {
		projectID = client.Store.Projects.GetIDByName(projectName)
		if projectID == "" {
			return fmt.Errorf("Did not find a project named '%v'", projectName)
		}
	}

	if projectID == "" {
		return fmt.Errorf("--project-id or --project-name flag is required")
	}

	if err := client.MoveSection(context.Background(), sectionID, projectID); err != nil {
		return err
	}

	return Sync(c)
}
