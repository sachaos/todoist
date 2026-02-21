package todoist

import (
	"testing"
)

func TestItem_LabelsString(t *testing.T) {
	store := &Store{
		Labels: Labels{
			{HaveID: HaveID{ID: "1"}, Name: "important"},
			{HaveID: HaveID{ID: "2"}, Name: "work"},
		},
	}
	
	item := Item{
		LabelNames: []string{"important", "work", "unknown_label"},
	}

	expected := "@important,@work,@unknown_label"
	result := item.LabelsString(store)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
