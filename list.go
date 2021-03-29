package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/acarl005/stripansi"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func traverseItems(item *todoist.Item, f func(item *todoist.Item, depth int), depth int) {
	f(item, depth)

	if item.ChildItem != nil {
		traverseItems(item.ChildItem, f, depth+1)
	}

	if item.BrotherItem != nil {
		traverseItems(item.BrotherItem, f, depth)
	}
}

type TaskJSON struct {
	ID       string `json:"id"`
	Priority string `json:"priority"`
	DueDate  string `json:"due_date"`
	Project  string `json:"project"`
	Labels   string `json:"labels"`
	Content  string `json:"content"`
}

func List(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	projectsCount := len(client.Store.Projects)
	projectIds := make([]int, projectsCount)
	for i, project := range client.Store.Projects {
		projectIds[i] = project.GetID()
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	rootItem := client.Store.RootItem

	if rootItem == nil {
		fmt.Fprintln(os.Stderr, "There is no task. You can fetch latest tasks by `todoist sync`.")
		return nil
	}

	isJson := c.GlobalBool("json")

	var jsonObjects []TaskJSON
	traverseItems(rootItem, func(item *todoist.Item, depth int) {
		r, err := Eval(ex, item, client.Store.Projects, client.Store.Labels)
		if err != nil {
			return
		}
		if !r || item.Checked == 1 {
			return
		}
		obj := TaskJSON{
			ID:       IdFormat(item),
			Priority: PriorityFormat(item.Priority),
			DueDate:  DueDateFormat(item.DateTime(), item.AllDay),
			Project: ProjectFormat(item.ProjectID, client.Store, projectColorHash, c) +
				SectionFormat(item.SectionID, client.Store, c),
			Labels:  item.LabelsString(client.Store),
			Content: ContentPrefix(client.Store, item, depth, c) + ContentFormat(item),
		}
		jsonObjects = append(jsonObjects, obj)
	}, 0)

	if c.Bool("priority") {
		// sort output by priority
		// and no need to use "else block" as items returned by API are already sorted by task id
		sort.Slice(jsonObjects, func(i, j int) bool {
			return stripansi.Strip(jsonObjects[i].Priority) < stripansi.Strip(jsonObjects[j].Priority)
		})
	}

	defer writer.Flush()

	if isJson {
		jsonData, err := json.Marshal(jsonObjects)
		if err != nil {
			return CommandFailed
		}
		writer.Write([]string{string(jsonData)})
	} else {
		if c.GlobalBool("header") {
			writer.Write([]string{"ID", "Priority", "DueDate", "Project", "Labels", "Content"})
		}

		for _, obj := range jsonObjects {
			writer.Write([]string{obj.ID, obj.Priority, obj.DueDate, obj.Project, obj.Labels, obj.Content})
		}
	}

	return nil
}
