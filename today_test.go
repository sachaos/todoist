package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/fatih/color"
	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func TestToday_NoArgs(t *testing.T) {
	color.NoColor = true

	store := testStore()
	client := todoist.NewClient(&todoist.Config{})
	client.Store = store

	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		"client": client,
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool("remote", false, "")
	flagSet.Bool("header", false, "")
	flagSet.Bool("color", false, "")
	flagSet.Bool("priority", false, "")
	flagSet.Bool("project-namespace", false, "")
	flagSet.Bool("indent", false, "")
	flagSet.Bool("namespace", false, "")
	flagSet.String("filter", "", "")
	flagSet.Int("limit", 0, "")
	flagSet.Parse([]string{})

	ctx := cli.NewContext(app, flagSet, nil)

	var buf bytes.Buffer
	writer = NewTSVWriter(&buf)

	// Should not error with no arguments
	err := Today(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
