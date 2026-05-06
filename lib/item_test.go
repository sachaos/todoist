package todoist

import (
	"testing"
)

func TestItem_AddParam_WithSectionID(t *testing.T) {
	item := Item{
		BaseItem: BaseItem{
			HaveProjectID: HaveProjectID{ProjectID: "proj-1"},
			Content:       "Test task",
		},
		HaveSectionID: HaveSectionID{SectionID: "sec-1"},
		Priority:      3,
	}
	param := item.AddParam().(map[string]interface{})

	if param["content"] != "Test task" {
		t.Errorf("expected 'Test task', got '%v'", param["content"])
	}
	if param["project_id"] != "proj-1" {
		t.Errorf("expected 'proj-1', got '%v'", param["project_id"])
	}
	if param["section_id"] != "sec-1" {
		t.Errorf("expected 'sec-1', got '%v'", param["section_id"])
	}
	if param["priority"] != 3 {
		t.Errorf("expected 3, got '%v'", param["priority"])
	}
}

func TestItem_AddParam_WithoutSectionID(t *testing.T) {
	item := Item{
		BaseItem: BaseItem{
			Content: "Test task",
		},
	}
	param := item.AddParam().(map[string]interface{})

	if _, ok := param["section_id"]; ok {
		t.Error("expected section_id to be absent when empty")
	}
}

func TestItem_MoveToSectionParam(t *testing.T) {
	item := &Item{
		BaseItem: BaseItem{
			HaveID: HaveID{ID: "item-123"},
		},
	}
	param := item.MoveToSectionParam("sec-456").(map[string]interface{})

	if param["id"] != "item-123" {
		t.Errorf("expected 'item-123', got '%v'", param["id"])
	}
	if param["section_id"] != "sec-456" {
		t.Errorf("expected 'sec-456', got '%v'", param["section_id"])
	}
}

func TestItem_UpdateParam_PriorityPreserved(t *testing.T) {
	item := Item{
		BaseItem: BaseItem{HaveID: HaveID{ID: "item-1"}},
		Priority: 3,
	}
	param := item.UpdateParam().(map[string]interface{})

	if param["priority"] != 3 {
		t.Errorf("expected priority 3 to be included, got %v", param["priority"])
	}
}

func TestItem_UpdateParam_ZeroPriorityOmitted(t *testing.T) {
	item := Item{
		BaseItem: BaseItem{HaveID: HaveID{ID: "item-1"}},
		Priority: 0,
	}
	param := item.UpdateParam().(map[string]interface{})

	if _, ok := param["priority"]; ok {
		t.Errorf("expected priority to be absent when zero, got %v", param["priority"])
	}
}

func TestItem_UpdateParam_ClearLabelNames(t *testing.T) {
	item := Item{
		BaseItem:        BaseItem{HaveID: HaveID{ID: "item-1"}},
		LabelNames:      []string{},
		ClearLabelNames: true,
	}
	param := item.UpdateParam().(map[string]interface{})

	v, ok := param["labels"]
	if !ok {
		t.Fatal("expected labels to be present when ClearLabelNames is true")
	}
	labels, ok := v.([]string)
	if !ok {
		t.Fatalf("expected labels to be []string, got %T", v)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty labels, got %v", labels)
	}
}

func TestItem_UpdateParam_EmptyLabelNamesOmittedWithoutClearFlag(t *testing.T) {
	item := Item{
		BaseItem:   BaseItem{HaveID: HaveID{ID: "item-1"}},
		LabelNames: []string{},
	}
	param := item.UpdateParam().(map[string]interface{})

	if _, ok := param["labels"]; ok {
		t.Error("expected labels to be absent when empty and ClearLabelNames is false")
	}
}

func TestItem_LabelsString(t *testing.T) {
	item1 := Item{
		LabelNames: []string{"important", "work", "unknown_label"},
	}

	expected1 := "@important,@work,@unknown_label"
	result1 := item1.LabelsString()
	if result1 != expected1 {
		t.Errorf("expected %s, got %s", expected1, result1)
	}

	item2 := Item{
		LabelNames: []string{},
	}
	expected2 := ""
	result2 := item2.LabelsString()
	if result2 != expected2 {
		t.Errorf("expected %q, got %q", expected2, result2)
	}
}
