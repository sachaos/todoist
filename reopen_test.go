package main

import (
	"testing"

	todoist "github.com/sachaos/todoist/lib"
)

func TestReopen_NoArgs(t *testing.T) {
	client := todoist.NewClient(&todoist.Config{})
	client.Store = testStore()

	ctx := newTestContext(client, []string{})

	err := Reopen(ctx)
	if err == nil {
		t.Error("expected error for no args, got nil")
	}
}
