package main

import (
	"regexp"
	"strconv"

	"github.com/sachaos/todoist/lib"
)

var priorityRegex = regexp.MustCompile("^p([1-4])$")

// Eval ...
func Eval(e Expression, item todoist.Item) (result bool, err error) {
	result = false
	matched := priorityRegex.FindStringSubmatch(e.(StringExpr).literal)

	if len(matched) == 0 {
		return
	} else {
		p, _ := strconv.Atoi(matched[1])
		if p == item.Priority {
			return true, err
		}
	}
	return
}
