package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

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

	itemList := [][]string{}
	project := client.Store.RootProject

	traverseProjects(project, func(pjt *todoist.Project, depth int) {
		itemList = append(itemList, []string{IdFormat(pjt), ProjectFormat(pjt.ID, client.Store, projectColorHash, c)})
	}, 0)

	defer writer.Flush()

	if c.GlobalBool("header") {
		writer.Write([]string{"ID", "Name"})
	}

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
