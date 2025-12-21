package main

import (
	"fmt"
	"os"
	"strings"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Cache(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	projectsCount := len(client.Store.Projects)
	projectIds := make([]string, projectsCount)
	for i, project := range client.Store.Projects {
		projectIds[i] = project.GetID()
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	itemList := [][]string{}
	syncedCount := 0
	unsyncedCount := 0

	rootItem := client.Store.RootItem

	if rootItem != nil {
		traverseItems(rootItem, func(item *todoist.Item, depth int) {
			r, err := Eval(ex, item, client.Store.Projects, client.Store.Labels)
			if err != nil {
				return
			}
			if !r || item.Checked {
				return
			}
			itemList = append(itemList, []string{
				IdFormat(item),
				PriorityFormat(item.Priority),
				DueDateFormat(item.DateTime(), item.AllDay),
				ProjectFormat(item.ProjectID, client.Store, projectColorHash, c) +
					SectionFormat(item.SectionID, client.Store, c),
				item.LabelsString(client.Store),
				ContentPrefix(client.Store, item, depth, c) + ContentFormat(item),
				"synced",
			})
			syncedCount++
		}, 0)
	}

	pc, err := GetPipelineCache(pipelineCachePath)
	if err == nil && !pc.IsEmpty() {
		pipelineItems := pc.GetItems()
		for _, pItem := range pipelineItems {
			if pItem.IsClose {
				continue
			}

			if pItem.IsQuick {
				itemList = append(itemList, []string{
					UnsyncedIdFormat("pending"),
					UnsyncedPriorityFormat(1), // Default priority
					UnsyncedDueDateFormat(""),
					UnsyncedProjectFormat(""),
					"",
					UnsyncedContentFormat(pItem.QuickText),
					UnsyncedStatusFormat("unsynced"),
				})
				unsyncedCount++
			} else {
				item := pItem.Item
				labelStr := ""
				if len(item.LabelNames) > 0 {
					labelStr = UnsyncedContentFormat("@" + strings.Join(item.LabelNames, ",@"))
				}
				itemList = append(itemList, []string{
					UnsyncedIdFormat(item.ID),
					UnsyncedPriorityFormat(item.Priority),
					UnsyncedDueDateFormat(dueDateString(item.DateTime(), item.AllDay)),
					UnsyncedProjectFormat(ProjectFormat(item.ProjectID, client.Store, projectColorHash, c)),
					labelStr,
					UnsyncedContentFormat(todoist.GetContentTitle(&item)),
					UnsyncedStatusFormat("unsynced"),
				})
				unsyncedCount++
			}
		}
	}

	if c.Bool("priority") == true {
		sortItems(&itemList, 1)
	}

	fmt.Fprintf(os.Stderr, "Total cached tasks: %d (synced: %d, unsynced: %d)\n",
		syncedCount+unsyncedCount, syncedCount, unsyncedCount)
	if unsyncedCount > 0 {
		fmt.Fprintf(os.Stderr, "Note: Unsynced tasks are displayed in yellow\n")
	}
	fmt.Fprintln(os.Stderr, "")
	defer writer.Flush()
	if c.Bool("header") {
		writer.Write([]string{"ID", "Priority", "DueDate", "Project", "Labels", "Content", "Status"})
	}

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
