package main

import (
	"sort"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func Filters(c *cli.Context) error {
	client := GetClient(c)

	defer writer.Flush()

	if c.Bool("header") {
		writer.Write([]string{"ID", "Name", "Query", "Favorite"})
	}

	// Sort filters by ItemOrder to preserve Todoist UI order
	type filterWithOrder struct {
		ID        string
		Name      string
		Query     string
		IsFavorite bool
		ItemOrder int
	}

	var filters []filterWithOrder
	for _, f := range client.Store.Filters {
		if f.IsDeleted {
			continue
		}
		filters = append(filters, filterWithOrder{
			ID:        f.ID,
			Name:      f.Name,
			Query:     f.Query,
			IsFavorite: f.IsFavorite,
			ItemOrder: f.ItemOrder,
		})
	}

	// Sort by ItemOrder
	sort.Slice(filters, func(i, j int) bool {
		return filters[i].ItemOrder < filters[j].ItemOrder
	})

	for _, f := range filters {
		favorite := "N"
		if f.IsFavorite {
			favorite = "Y"
		}
		writer.Write([]string{color.BlueString(f.ID), f.Name, f.Query, favorite})
	}

	return nil
}
