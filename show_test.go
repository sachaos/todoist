package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/fatih/color"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func newTestContext(client *todoist.Client, args []string) *cli.Context {
	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		"client": client,
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool("browse", false, "")
	flagSet.Bool("header", false, "")
	flagSet.Bool("color", false, "")
	flagSet.Bool("project-namespace", false, "")
	flagSet.Bool("indent", false, "")
	flagSet.Parse(args)

	ctx := cli.NewContext(app, flagSet, nil)
	return ctx
}

func testStore() *todoist.Store {
	store := &todoist.Store{
		Projects: todoist.Projects{
			{HaveID: todoist.HaveID{ID: "proj-1"}, Name: "Work"},
		},
		Sections: todoist.Sections{
			{
				HaveID:        todoist.HaveID{ID: "sec-1"},
				HaveProjectID: todoist.HaveProjectID{ProjectID: "proj-1"},
				Name:          "Backlog",
			},
		},
		Items: todoist.Items{
			{
				BaseItem: todoist.BaseItem{
					HaveID:        todoist.HaveID{ID: "item-1"},
					HaveProjectID: todoist.HaveProjectID{ProjectID: "proj-1"},
					Content:       "Test task",
				},
				HaveSectionID: todoist.HaveSectionID{SectionID: "sec-1"},
				Priority:      3,
			},
		},
	}
	store.ConstructItemTree()
	return store
}

func TestShow_Item(t *testing.T) {
	color.NoColor = true

	store := testStore()
	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	ctx := newTestContext(client, []string{"item-1"})

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	err := Show(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("item-1")) {
		t.Errorf("expected output to contain item ID, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Test task")) {
		t.Errorf("expected output to contain item content, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Backlog")) {
		t.Errorf("expected output to contain section name, got: %s", output)
	}
}

func TestShow_Section(t *testing.T) {
	color.NoColor = true

	store := testStore()
	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	ctx := newTestContext(client, []string{"sec-1"})

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	err := Show(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("sec-1")) {
		t.Errorf("expected output to contain section ID, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Backlog")) {
		t.Errorf("expected output to contain section name, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Work")) {
		t.Errorf("expected output to contain project name, got: %s", output)
	}
}

func TestShow_NotFound(t *testing.T) {
	color.NoColor = true

	store := testStore()
	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	ctx := newTestContext(client, []string{"nonexistent"})

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	err := Show(ctx)
	if err != IdNotFound {
		t.Errorf("expected IdNotFound error, got %v", err)
	}
}
