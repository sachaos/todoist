package main

import (
	"testing"

	"github.com/sachaos/todoist/lib"
	"github.com/stretchr/testify/assert"
)

func testFilterEval(t *testing.T, f string, item todoist.Item, expect bool) {
	actual, _ := Eval(Filter(f), item)
	assert.Equal(t, expect, actual, "they should be equal")
}

func TestEval(t *testing.T) {
	testFilterEval(t, "p1", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p2", todoist.Item{Priority: 1}, false)

	testFilterEval(t, "", todoist.Item{}, true)

	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 2}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 3}, false)

	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 3}, false)
}
