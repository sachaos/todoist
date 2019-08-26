package main

import (
	"fmt"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"os"
)

func traverseItems(item *todoist.Item, f func(item *todoist.Item, depth int), depth int) {
	f(item, depth)

	if item.ChildItem != nil {
		traverseItems(item.ChildItem, f, depth + 1)
	}

	if item.BrotherItem != nil {
		traverseItems(item.BrotherItem, f, depth)
	}
}

func List(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	projectsCount := len(client.Store.Projects)
	projectIds := make([]int, projectsCount, projectsCount)
	for i, project := range client.Store.Projects {
		projectIds[i] = project.GetID()
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	itemList := [][]string{}
	rootItem := client.Store.RootItem

	if rootItem == nil {
		fmt.Fprintln(os.Stderr, "There is no task. You can fetch latest tasks by `todoist sync`.")
		return nil
	}

	traverseItems(rootItem, func(item *todoist.Item, depth int) {
		r, err := Eval(ex, item, client.Store.Projects, client.Store.Labels)
		if err != nil {
			return
		}
		if !r || item.Checked == 1 {
			return
		}
		itemList = append(itemList, []string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DateTime(), item.AllDay),
			ProjectFormat(item.ProjectID, client.Store, projectColorHash, c),
			item.LabelsString(client.Store),
			ContentPrefix(client.Store, item, depth, c) + ContentFormat(item),
		})
	}, 0)

	defer writer.Flush()

	if c.GlobalBool("header") {
		writer.Write([]string{"ID", "Priority", "DueDate", "Project", "Labels", "Content"})
	}

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
