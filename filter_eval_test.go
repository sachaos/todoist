package main

import (
	"testing"

	"github.com/sachaos/todoist/lib"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	r, _ := Eval(Filter("p1"), todoist.Item{Priority: 1})
	assert.Equal(t, r, true, "they should be equal")
	r, _ = Eval(Filter("p2"), todoist.Item{Priority: 1})
	assert.Equal(t, r, false, "they should be equal")
}
