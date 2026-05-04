package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/fatih/color"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func newFilterTestContext(client *todoist.Client) *cli.Context {
	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		"client": client,
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool("header", false, "")
	flagSet.Bool("color", false, "")
	flagSet.Parse([]string{})

	ctx := cli.NewContext(app, flagSet, nil)
	return ctx
}

func TestFilters_NoArgs(t *testing.T) {
	color.NoColor = true

	store := &todoist.Store{
		Filters: []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Query     string `json:"query"`
			Color     string `json:"color"`
			ItemOrder int    `json:"item_order"`
			IsDeleted bool   `json:"is_deleted"`
			IsFavorite bool  `json:"is_favorite"`
			IsFrozen  bool   `json:"is_frozen"`
		}{
			{
				ID:        "filter-1",
				Name:      "My Overdue",
				Query:     "overdue & !#Work",
				Color:     "berry_red",
				ItemOrder: 1,
				IsDeleted: false,
				IsFavorite: true,
				IsFrozen:  false,
			},
			{
				ID:        "filter-2",
				Name:      "This Week",
				Query:     "(today | next 7 days) & !#Someday",
				Color:     "blue",
				ItemOrder: 2,
				IsDeleted: false,
				IsFavorite: false,
				IsFrozen:  false,
			},
			{
				ID:        "filter-3",
				Name:      "Deleted Filter",
				Query:     "overdue",
				Color:     "red",
				ItemOrder: 3,
				IsDeleted: true,
				IsFavorite: false,
				IsFrozen:  false,
			},
		},
	}

	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	ctx := newFilterTestContext(client)

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	err := Filters(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()

	// Verify deleted filter is skipped
	if bytes.Contains([]byte(output), []byte("filter-3")) {
		t.Errorf("expected output to skip deleted filter, got: %s", output)
	}

	// Verify non-deleted filters are included
	if !bytes.Contains([]byte(output), []byte("filter-1")) {
		t.Errorf("expected output to contain filter-1, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("filter-2")) {
		t.Errorf("expected output to contain filter-2, got: %s", output)
	}

	// Verify names are included
	if !bytes.Contains([]byte(output), []byte("My Overdue")) {
		t.Errorf("expected output to contain 'My Overdue', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("This Week")) {
		t.Errorf("expected output to contain 'This Week', got: %s", output)
	}

	// Verify queries are included
	if !bytes.Contains([]byte(output), []byte("overdue & !#Work")) {
		t.Errorf("expected output to contain 'overdue & !#Work', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("(today | next 7 days) & !#Someday")) {
		t.Errorf("expected output to contain '(today | next 7 days) & !#Someday', got: %s", output)
	}

	// Verify favorite status is included
	if !bytes.Contains([]byte(output), []byte("Y")) {
		t.Errorf("expected output to contain 'Y' for favorite, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("N")) {
		t.Errorf("expected output to contain 'N' for not favorite, got: %s", output)
	}
}

func TestFilters_WithHeader(t *testing.T) {
	color.NoColor = true

	store := &todoist.Store{
		Filters: []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Query     string `json:"query"`
			Color     string `json:"color"`
			ItemOrder int    `json:"item_order"`
			IsDeleted bool   `json:"is_deleted"`
			IsFavorite bool  `json:"is_favorite"`
			IsFrozen  bool   `json:"is_frozen"`
		}{
			{
				ID:        "filter-1",
				Name:      "Test Filter",
				Query:     "overdue",
				Color:     "red",
				ItemOrder: 1,
				IsDeleted: false,
				IsFavorite: false,
				IsFrozen:  false,
			},
		},
	}

	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		"client": client,
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool("header", true, "")
	flagSet.Bool("color", false, "")
	flagSet.Parse([]string{})

	ctx := cli.NewContext(app, flagSet, nil)

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	err := Filters(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()

	// Verify header is included
	if !bytes.Contains([]byte(output), []byte("ID")) {
		t.Errorf("expected output to contain header 'ID', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Name")) {
		t.Errorf("expected output to contain header 'Name', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Query")) {
		t.Errorf("expected output to contain header 'Query', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Favorite")) {
		t.Errorf("expected output to contain header 'Favorite', got: %s", output)
	}
}
