package main

import (
	"context"
	"encoding/json"

	todoist "github.com/sachaos/todoist/lib"

	"github.com/urfave/cli"
)

type CompletedJSON struct {
	ID            string `json:"id"`
	CompletedDate string `json:"completed_date"`
	Project       string `json:"project"`
	Content       string `json:"content"`
}

func CompletedList(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []int
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	var completed todoist.Completed

	if err := client.CompletedAll(context.Background(), &completed); err != nil {
		return err
	}

	isJson := c.GlobalBool("json")

	defer writer.Flush()

	if !isJson && c.GlobalBool("header") {
		writer.Write([]string{"ID", "CompletedDate", "Project", "Content"})
	}

	var jsonObjects []CompletedJSON
	for _, item := range completed.Items {
		result, err := Eval(ex, item, client.Store.Projects, client.Store.Labels)
		if err != nil {
			return err
		}
		if !result {
			continue
		}

		obj := CompletedJSON{
			ID:            IdFormat(item),
			CompletedDate: CompletedDateFormat(item.DateTime()),
			Project:       ProjectFormat(item.ProjectID, client.Store, projectColorHash, c),
			Content:       ContentFormat(item),
		}
		if isJson {
			jsonObjects = append(jsonObjects, obj)
		} else {
			writer.Write([]string{obj.ID, obj.CompletedDate, obj.Project, obj.Content})
		}

	}

	if isJson {
		jsonData, err := json.Marshal(jsonObjects)
		if err != nil {
			return CommandFailed
		}
		writer.Write([]string{string(jsonData)})
	}

	return nil
}
