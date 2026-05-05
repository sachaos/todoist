package main

import (
	"errors"
	"testing"

	todoist "github.com/sachaos/todoist/lib"
)

func TestModify_NoArgs(t *testing.T) {
	client := todoist.NewClient(&todoist.Config{})
	client.Store = testStore()

	ctx := newTestContext(client, []string{})

	err := Modify(ctx)
	if err == nil {
		t.Error("expected error for no args, got nil")
	}
}

func TestModify_NotFound(t *testing.T) {
	client := todoist.NewClient(&todoist.Config{})
	client.Store = testStore()

	ctx := newTestContext(client, []string{"nonexistent-id"})

	err := Modify(ctx)
	if !errors.Is(err, IdNotFound) {
		t.Errorf("expected IdNotFound, got %v", err)
	}
}
