package main

import (
	"container/list"
	"sort"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func List(sync todoist.Sync, c *cli.Context) error {
	colorList := ColorList()
	projectTree := ProjectTree(sync)
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	itemQue := list.New()
	sort.Sort(sync.Items)
	for _, item := range sync.Items {
		itemQue.PushBack(item)
	}

	tree, err := NewTree(itemQue)
	if err != nil {
		return err
	}

	defer writer.Flush()

	for _, node := range tree.Traverse() {
		item := node.Value.(todoist.Item)
		if item.Checked == 1 {
			continue
		}
		writer.Write([]string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DueDateTime(), item.AllDay),
			ProjectFormat(item.ProjectID, projectTree, projectColorHash, c),
			item.LabelsString(sync.Labels),
			ContentPrefix(node, c) + ContentFormat(item),
		})
	}

	return nil
}
