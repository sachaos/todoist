package main

import (
	"encoding/json"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

type LabelJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Labels(c *cli.Context) error {
	client := GetClient(c)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	isJson := c.GlobalBool("json")

	defer writer.Flush()
	if !isJson && c.GlobalBool("header") {
		writer.Write([]string{"ID", "Name"})
	}

	var jsonObjects []LabelJSON
	for _, label := range client.Store.Labels {
		obj := LabelJSON{
			ID:   IdFormat(label),
			Name: "@" + label.Name,
		}
		if isJson {
			jsonObjects = append(jsonObjects, obj)
		} else {
			writer.Write([]string{obj.ID, obj.Name})
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
