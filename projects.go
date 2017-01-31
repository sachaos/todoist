package main

import (
	"container/list"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func ProjectTree(sync todoist.Sync) *Tree {
	itemQue := list.New()
	for _, item := range sync.Projects {
		itemQue.PushBack(item)
	}

	tree, _ := NewTree(itemQue)
	return tree
}

func Projects(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectTree := ProjectTree(sync)
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	defer writer.Flush()

	for _, node := range projectTree.Traverse() {
		project := node.Value.(todoist.Project)
		writer.Write([]string{IdFormat(project), ProjectFormat(project.ID, projectTree, projectColorHash, c)})
	}

	return nil
}
