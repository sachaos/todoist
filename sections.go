package main

import (
	"github.com/urfave/cli/v2"
)

func Sections(c *cli.Context) error {
	client := GetClient(c)

	defer writer.Flush()

	if c.Bool("header") {
		writer.Write([]string{"ID", "Project", "Name"})
	}

	for _, section := range client.Store.Sections {
		project := client.Store.FindProject(section.ProjectID)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}
		writer.Write([]string{IdFormat(section), projectName, section.Name})
	}

	return nil
}
