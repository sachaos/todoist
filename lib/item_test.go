package todoist

import (
	"testing"
)

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
