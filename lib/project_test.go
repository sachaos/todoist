package todoist

import (
	"testing"
)

func TestProjects_GetIDByName(t *testing.T) {
	projects := Projects{
		Project{HaveID: HaveID{ID: "1"}, Name: "Inbox"},
		Project{HaveID: HaveID{ID: "2"}, Name: "Packing"},
		Project{HaveID: HaveID{ID: "3"}, Name: "Work"},
	}

	// Bare name lookup works
	if id := projects.GetIDByName("Packing"); id != "2" {
		t.Errorf("expected '2', got '%s'", id)
	}

	// # prefix is stripped before matching, matching ProjectFormat display output
	if id := projects.GetIDByName("#Packing"); id != "2" {
		t.Errorf("expected '2', got '%s'", id)
	}

	// Not found returns empty string
	if id := projects.GetIDByName("Nonexistent"); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}

	// # prefix alone (no name) returns empty string
	if id := projects.GetIDByName("#"); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}

	// Empty projects returns empty string
	empty := Projects{}
	if id := empty.GetIDByName("Packing"); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}
}
