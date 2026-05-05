package main

import (
	"testing"

	todoist "github.com/sachaos/todoist/lib"
)

func TestVersion_DevBuild(t *testing.T) {
	orig := version
	version = ""
	defer func() { version = orig }()

	ctx := newTestContext(todoist.NewClient(&todoist.Config{}), []string{})
	if err := Version(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVersion_WithCommitAndDate(t *testing.T) {
	origVersion, origCommit, origDate := version, commit, date
	version = "0.21.0"
	commit = "abc1234"
	date = "2026-05-04"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	ctx := newTestContext(todoist.NewClient(&todoist.Config{}), []string{})
	if err := Version(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVersion_VersionOnly(t *testing.T) {
	origVersion, origCommit, origDate := version, commit, date
	version = "0.21.0"
	commit = ""
	date = ""
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	ctx := newTestContext(todoist.NewClient(&todoist.Config{}), []string{})
	if err := Version(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
