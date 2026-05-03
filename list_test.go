package main

import (
	"flag"
	"testing"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func TestListRemote_NoFilter(t *testing.T) {
	client := todoist.NewClient(&todoist.Config{})
	client.Store = testStore()

	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		"client": client,
	}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool("remote", false, "")
	flagSet.String("filter", "", "")
	flagSet.Int("limit", 0, "")
	flagSet.Parse([]string{})
	flagSet.Set("remote", "true")

	ctx := cli.NewContext(app, flagSet, nil)

	err := List(ctx)
	if err == nil {
		t.Fatal("expected error when --remote used without --filter, got nil")
	}
	if err.Error() != "--remote requires --filter" {
		t.Errorf("expected '--remote requires --filter', got: %v", err)
	}
}
