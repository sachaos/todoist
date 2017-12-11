package main

import (
	"testing"
	"time"

	"github.com/sachaos/todoist/lib"
	"github.com/stretchr/testify/assert"
)

var testTimeZone = time.FixedZone("JST", 9*60*60)

func testFilterEval(t *testing.T, f string, item todoist.Item, expect bool) {
	actual, _ := Eval(Filter(f), item)
	assert.Equal(t, expect, actual, "they should be equal")
}

func TestEval(t *testing.T) {
	testFilterEval(t, "", todoist.Item{}, true)
}

func TestPriorityEval(t *testing.T) {
	testFilterEval(t, "p1", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p2", todoist.Item{Priority: 1}, false)
}

func TestBoolInfixOpExp(t *testing.T) {
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 2}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 3}, false)

	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 3}, false)
}

func TestNotOpEval(t *testing.T) {
	testFilterEval(t, "!p1", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "!(p1 | p2)", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "!(p1 | p2)", todoist.Item{Priority: 3}, true)
}

func TestDueOnEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 15:00:00 +0000"}, true)  // JST: Mon 2 Oct 2017 00:00:00
	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 14:59:59 +0000"}, true)  // JST: Mon 2 Oct 2017 23:59:59
	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 15:00:00 +0000"}, false) // JST: Tue 3 Oct 2017 00:00:00

	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 14:59:59 +0000"}, true)   // JST: Sun 1 Oct 2017 23:59:59
	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sat 30 Sep 2017 15:00:00 +0000"}, true)  // JST: Sun 1 Oct 2017 00:00:00
	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sat 30 Sep 2017 14:59:59 +0000"}, false) // JST: Sat 30 Sept 2017 23:59:59
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 15:00:00 +0000"}, true)    // JST: Tue 3 Oct 2017 00:00:00
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Tue 3 Oct 2017 14:59:59 +0000"}, true)    // JST: Tue 3 Oct 2017 23:59:59
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Tue 3 Oct 2017 15:00:00 +0000"}, false)   // JST: Wed 4 Oct 2017 00:00:00

	testFilterEval(t, "10/2/2017", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 01:00:00 +0000"}, true)        // JST: Mon 2 Oct 2017 10:00:00
	testFilterEval(t, "10/2/2017 10:00", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 01:00:00 +0000"}, false) // JST: Mon 2 Oct 2017 10:00:00
}

func TestNoDateEval(t *testing.T) {
	testFilterEval(t, "no date", todoist.Item{DueDateUtc: ""}, true)
	testFilterEval(t, "no due date", todoist.Item{DueDateUtc: ""}, true)

	testFilterEval(t, "no date", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 15:00:00 +0000"}, false) // JST: Mon 2 Oct 2017 00:00:00
}

func TestDueBeforeEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	testFilterEval(t, "due before: 10/2/2017", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 15:00:00 +0000"}, false)      // JST: Mon 2 Oct 2017 00:00:00
	testFilterEval(t, "due before: 10/2/2017", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 14:59:59 +0000"}, true)       // JST: Sun 1 Oct 2017 23:59:59
	testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 4:00:00 +0000"}, false) // JST: Mon 2 Oct 2017 13:00:00
	testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 3:59:00 +0000"}, true)  // JST: Mon 2 Oct 2017 12:59:00

	testFilterEval(t, "due before: 10/2/2017 13:00", todoist.Item{DueDateUtc: ""}, false) // JST: Mon 2 Oct 2017 12:59:00
}

func TestOverDueEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 12, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 12:00:00
	setNow(timeNow)

	testFilterEval(t, "over due", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 2:59:00 +0000"}, true)  // JST: Mon 2 Oct 2017 11:59:00
	testFilterEval(t, "over due", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 3:00:00 +0000"}, false) // JST: Mon 2 Oct 2017 12:00:00
	testFilterEval(t, "od", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 2:59:00 +0000"}, true)        // JST: Mon 2 Oct 2017 11:59:00
	testFilterEval(t, "od", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 3:00:00 +0000"}, false)       // JST: Mon 2 Oct 2017 12:00:00

	testFilterEval(t, "od", todoist.Item{DueDateUtc: ""}, false) // JST: Mon 2 Oct 2017 12:00:00
}

func TestDueAfterEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, testTimeZone) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	testFilterEval(t, "due after: 10/2/2017", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 14:59:59 +0000"}, false)      // JST: Mon 2 Oct 2017 23:59:59
	testFilterEval(t, "due after: 10/2/2017", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 15:00:00 +0000"}, true)       // JST: Tue 3 Oct 2017 00:00:00
	testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 4:00:00 +0000"}, false) // JST: Mon 2 Oct 2017 13:00:00
	testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 4:01:00 +0000"}, true)  // JST: Mon 2 Oct 2017 13:01:00

	testFilterEval(t, "due after: 10/2/2017 13:00", todoist.Item{DueDateUtc: ""}, false) // JST: Mon 2 Oct 2017 13:01:00
}
