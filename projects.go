package main

import (
	"encoding/json"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

type ProjectJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func traverseProjects(pjt *todoist.Project, f func(pjt *todoist.Project, depth int), depth int) {
	f(pjt, depth)

	if pjt.ChildProject != nil {
		traverseProjects(pjt.ChildProject, f, depth+1)
	}

	if pjt.BrotherProject != nil {
		traverseProjects(pjt.BrotherProject, f, depth)
	}
}

func Projects(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []string
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	project := client.Store.RootProject

	var jsonObjects []ProjectJSON
	traverseProjects(project, func(pjt *todoist.Project, depth int) {
		jsonObjects = append(jsonObjects, ProjectJSON{
			ID:   IdFormat(pjt),
			Name: ProjectFormat(pjt.ID, client.Store, projectColorHash, c),
		})
	}, 0)

	defer writer.Flush()

	if c.GlobalBool("json") {
		jsonData, err := json.Marshal(jsonObjects)
		if err != nil {
			return CommandFailed
		}
		writer.Write([]string{string(jsonData)})
	} else {
		if c.GlobalBool("header") {
			writer.Write([]string{"ID", "Name"})
		}

		for _, obj := range jsonObjects {
			writer.Write([]string{obj.ID, obj.Name})
		}
	}

	return nil
}
