package todoist

import (
	"testing"
)

func TestSections_Active(t *testing.T) {
	sections := Sections{
		Section{HaveID: HaveID{ID: "1"}, Name: "Active"},
		Section{HaveID: HaveID{ID: "2"}, Name: "Deleted", IsDeleted: true},
		Section{HaveID: HaveID{ID: "3"}, Name: "Archived", IsArchived: true},
		Section{HaveID: HaveID{ID: "4"}, Name: "Also Active"},
	}

	active := sections.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active sections, got %d", len(active))
	}
	if active[0].ID != "1" || active[1].ID != "4" {
		t.Errorf("unexpected active sections: %v", active)
	}
}

func TestSections_GetIDByName(t *testing.T) {
	sections := Sections{
		Section{
			HaveID:        HaveID{ID: "100"},
			HaveProjectID: HaveProjectID{ProjectID: "1"},
			Name:          "Backlog",
		},
		Section{
			HaveID:        HaveID{ID: "200"},
			HaveProjectID: HaveProjectID{ProjectID: "1"},
			Name:          "In Progress",
		},
		Section{
			HaveID:        HaveID{ID: "300"},
			HaveProjectID: HaveProjectID{ProjectID: "2"},
			Name:          "Backlog",
		},
		Section{
			HaveID:        HaveID{ID: "400"},
			HaveProjectID: HaveProjectID{ProjectID: "2"},
			Name:          "Archived Backlog",
			IsArchived:    true,
		},
		Section{
			HaveID:        HaveID{ID: "500"},
			HaveProjectID: HaveProjectID{ProjectID: "2"},
			Name:          "Deleted Backlog",
			IsDeleted:     true,
		},
	}

	// Project-scoped lookup returns the correct section
	if id := sections.GetIDByName("Backlog", "1"); id != "100" {
		t.Errorf("expected '100', got '%s'", id)
	}
	if id := sections.GetIDByName("Backlog", "2"); id != "300" {
		t.Errorf("expected '300', got '%s'", id)
	}

	// Global lookup (no project) returns first match
	if id := sections.GetIDByName("In Progress", ""); id != "200" {
		t.Errorf("expected '200', got '%s'", id)
	}

	// Archived and deleted sections are never returned
	if id := sections.GetIDByName("Archived Backlog", ""); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}
	if id := sections.GetIDByName("Deleted Backlog", ""); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}

	// Not found
	if id := sections.GetIDByName("Nonexistent", ""); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}

	// Empty sections
	empty := Sections{}
	if id := empty.GetIDByName("Backlog", ""); id != "" {
		t.Errorf("expected '', got '%s'", id)
	}
}

func TestSection_AddParam(t *testing.T) {
	// Both name and project ID
	section := Section{
		HaveProjectID: HaveProjectID{ProjectID: "proj-1"},
		Name:          "My Section",
	}
	param := section.AddParam().(map[string]interface{})
	if param["name"] != "My Section" {
		t.Errorf("expected 'My Section', got '%v'", param["name"])
	}
	if param["project_id"] != "proj-1" {
		t.Errorf("expected 'proj-1', got '%v'", param["project_id"])
	}

	// Name only, no project ID
	section2 := Section{
		Name: "Another Section",
	}
	param2 := section2.AddParam().(map[string]interface{})
	if param2["name"] != "Another Section" {
		t.Errorf("expected 'Another Section', got '%v'", param2["name"])
	}
	if _, ok := param2["project_id"]; ok {
		t.Error("expected project_id to be absent")
	}

	// Empty section
	section3 := Section{}
	param3 := section3.AddParam().(map[string]interface{})
	if len(param3) != 0 {
		t.Errorf("expected empty param, got %v", param3)
	}
}

func TestStore_FindSection(t *testing.T) {
	store := Store{
		Sections: Sections{
			Section{HaveID: HaveID{ID: "s1"}, Name: "Section 1"},
			Section{HaveID: HaveID{ID: "s2"}, Name: "Section 2"},
		},
	}
	store.SectionMap = map[string]*Section{}
	for i := range store.Sections {
		store.SectionMap[store.Sections[i].ID] = &store.Sections[i]
	}

	found := store.FindSection("s1")
	if found == nil {
		t.Fatal("expected to find section s1")
	}
	if found.Name != "Section 1" {
		t.Errorf("expected 'Section 1', got '%s'", found.Name)
	}

	found2 := store.FindSection("s2")
	if found2 == nil {
		t.Fatal("expected to find section s2")
	}
	if found2.Name != "Section 2" {
		t.Errorf("expected 'Section 2', got '%s'", found2.Name)
	}

	notFound := store.FindSection("nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent section")
	}
}
