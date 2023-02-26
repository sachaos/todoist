package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sachaos/todoist/lib"
)

const DateFormat = "Mon 2 Jan 2006 15:04:05 +0000"

var testTimeZone = time.FixedZone("JST", 9*60*60)

func due(s string) *todoist.Due {
	t, _ := time.Parse(DateFormat, s)
	t = t.In(testTimeZone)
	date := t.Format(todoist.RFC3339DateTime)
	return &todoist.Due{
		Date: date,
	}
}

func testFilterEval(t *testing.T, f string, item todoist.Item, expect bool) {
	actual, _ := Eval(Filter(f), &item, todoist.Projects{}, todoist.Labels{})
	assert.Equal(t, expect, actual, "they should be equal")
}

func testFilterEvalWithProject(t *testing.T, f string, item todoist.Item, projects todoist.Projects, expect bool) {
	actual, _ := Eval(Filter(f), &item, projects, todoist.Labels{})
	assert.Equal(t, expect, actual, "they should be equal")
}

func testFilterEvalWithLabel(t *testing.T, f string, item todoist.Item, labels todoist.Labels, expect bool) {
	actual, _ := Eval(Filter(f), &item, todoist.Projects{}, labels)
	assert.Equal(t, expect, actual, "they should be equal")
}

func TestEval(t *testing.T) {
	testFilterEval(t, "", todoist.Item{}, true)
}

func TestPriorityEval(t *testing.T) {
	testFilterEval(t, "p4", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p3", todoist.Item{Priority: 1}, false)
}

func TestLabelEval(t *testing.T) {
	labels := todoist.Labels{
		todoist.Label{
			HaveID: todoist.HaveID{ID: "1"},
			Name:   "must",
		},
		todoist.Label{
			HaveID: todoist.HaveID{ID: "2"},
			Name:   "icebox",
		}, todoist.Label{
			HaveID: todoist.HaveID{ID: "3"},
			Name:   "another",
		},
	}

	item1 := todoist.Item{}
	item1.LabelNames = []string{"1", "2"}

	// testFilterEvalWithLabel(t, "@must", item1, labels, true)
	// testFilterEvalWithLabel(t, "@icebox", item1, labels, true)
	testFilterEvalWithLabel(t, "@another", item1, labels, false)
}

func TestProjectEval(t *testing.T) {
	projects := todoist.Projects{
		todoist.Project{
			HaveID: todoist.HaveID{ID: "1"},
			Name:   "private",
		},
		todoist.Project{
			HaveID:       todoist.HaveID{ID: "2"},
			HaveParentID: todoist.HaveParentID{ParentID: &[]string{"1"}[0]},
			Name:         "nested",
		},
	}

	item1 := todoist.Item{}
	item1.ProjectID = "1"

	item2 := todoist.Item{}
	item2.ProjectID = "2"

	testFilterEvalWithProject(t, "#private", item1, projects, true)
	testFilterEvalWithProject(t, "#hoge", item1, projects, false)
	testFilterEvalWithProject(t, "#private", item2, projects, false)
	testFilterEvalWithProject(t, "##private", item2, projects, true)
}

func TestBoolInfixOpExp(t *testing.T) {
	testFilterEval(t, "p3 | p4", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p3 | p4", todoist.Item{Priority: 2}, true)
	testFilterEval(t, "p3 | p4", todoist.Item{Priority: 3}, false)

	testFilterEval(t, "p3 & p4", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "p3 & p4", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "p3 & p4", todoist.Item{Priority: 3}, false)
}

func TestNotOpEval(t *testing.T) {
	testFilterEval(t, "!p4", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "!(p3 | p4)", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "!(p3 | p4)", todoist.Item{Priority: 3}, true)
}

func TestDueOnEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	// testFilterEval(t, "today", todoist.Item{Due: due("Sun 1 Oct 2017 15:00:00 +0000")}, true)  // JST: Mon 2 Oct 2017 00:00:00
	// testFilterEval(t, "today", todoist.Item{Due: due("Mon 2 Oct 2017 14:59:59 +0000")}, true)  // JST: Mon 2 Oct 2017 23:59:59
	testFilterEval(t, "today", todoist.Item{Due: due("Mon 2 Oct 2017 15:00:00 +0000")}, false) // JST: Tue 3 Oct 2017 00:00:00

	// testFilterEval(t, "yesterday", todoist.Item{Due: due("Sun 1 Oct 2017 14:59:59 +0000")}, true)   // JST: Sun 1 Oct 2017 23:59:59
	// testFilterEval(t, "yesterday", todoist.Item{Due: due("Sat 30 Sep 2017 15:00:00 +0000")}, true)  // JST: Sun 1 Oct 2017 00:00:00
	// testFilterEval(t, "yesterday", todoist.Item{Due: due("Sat 30 Sep 2017 14:59:59 +0000")}, false) // JST: Sat 30 Sept 2017 23:59:59
	// testFilterEval(t, "tomorrow", todoist.Item{Due: due("Mon 2 Oct 2017 15:00:00 +0000")}, true)  // JST: Tue 3 Oct 2017 00:00:00
	// testFilterEval(t, "tomorrow", todoist.Item{Due: due("Tue 3 Oct 2017 14:59:59 +0000")}, true)  // JST: Tue 3 Oct 2017 23:59:59
	testFilterEval(t, "tomorrow", todoist.Item{Due: due("Tue 3 Oct 2017 15:00:00 +0000")}, false) // JST: Wed 4 Oct 2017 00:00:00

	testFilterEval(t, "10/2/2017", todoist.Item{Due: due("Mon 2 Oct 2017 01:00:00 +0000")}, true)        // JST: Mon 2 Oct 2017 10:00:00
	testFilterEval(t, "10/2/2017 10:00", todoist.Item{Due: due("Mon 2 Oct 2017 01:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 10:00:00
}

func TestNoDateEval(t *testing.T) {
	testFilterEval(t, "no date", todoist.Item{Due: nil}, true)
	testFilterEval(t, "no due date", todoist.Item{Due: nil}, true)

	testFilterEval(t, "no date", todoist.Item{Due: due("Sun 1 Oct 2017 15:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 00:00:00
}

func TestDueBeforeEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	testFilterEval(t, "due before: 10/2/2017", todoist.Item{Due: due("Sun 1 Oct 2017 15:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 00:00:00
	// testFilterEval(t, "due before: 10/2/2017", todoist.Item{Due: due("Sun 1 Oct 2017 14:59:59 +0000")}, true)  // JST: Sun 1 Oct 2017 23:59:59
	testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{Due: due("Mon 2 Oct 2017 4:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 13:00:00
	// testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{Due: due("Mon 2 Oct 2017 3:59:00 +0000")}, true)  // JST: Mon 2 Oct 2017 12:59:00

	testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{Due: nil}, false) // JST: Mon 2 Oct 2017 12:59:00
}

func TestOverDueEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 12, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 12:00:00
	setNow(timeNow)

	// testFilterEval(t, "over due", todoist.Item{Due: due("Mon 2 Oct 2017 2:59:00 +0000")}, true) // JST: Mon 2 Oct 2017 11:59:00
	// testFilterEval(t, "over due", todoist.Item{Due: due("Mon 2 Oct 2017 3:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 12:00:00
	// testFilterEval(t, "od", todoist.Item{Due: due("Mon 2 Oct 2017 2:59:00 +0000")}, true)        // JST: Mon 2 Oct 2017 11:59:00
	// testFilterEval(t, "od", todoist.Item{Due: due("Mon 2 Oct 2017 3:00:00 +0000")}, false)       // JST: Mon 2 Oct 2017 12:00:00

	// testFilterEval(t, "od", todoist.Item{Due: nil}, false) // JST: Mon 2 Oct 2017 12:00:00
}

func TestDueAfterEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	//testFilterEval(t, "due after: 10/2/2017", todoist.Item{Due: due("Mon 2 Oct 2017 14:59:59 +0000")}, false)      // JST: Mon 2 Oct 2017 23:59:59
	//testFilterEval(t, "due after: 10/2/2017", todoist.Item{Due: due("Mon 2 Oct 2017 15:00:00 +0000")}, true)       // JST: Tue 3 Oct 2017 00:00:00
	//testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{Due: due("Mon 2 Oct 2017 4:00:00 +0000")}, false) // JST: Mon 2 Oct 2017 13:00:00
	//testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{Due: due("Mon 2 Oct 2017 4:01:00 +0000")}, true)  // JST: Mon 2 Oct 2017 13:01:00
	//
	//testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{Due: nil}, false) // JST: Mon 2 Oct 2017 13:01:00
}
