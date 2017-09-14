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
	switch e.(type) {
	case BoolInfixOpExpr:
		e := e.(BoolInfixOpExpr)
		lr, err := Eval(e.left, item)
		rr, err := Eval(e.right, item)
		if err != nil {
			return false, nil
		}
		switch e.operator {
		case '&':
			return lr && rr, nil
		case '|':
			return lr || rr, nil
		}
	case StringExpr:
		matched := priorityRegex.FindStringSubmatch(e.(StringExpr).literal)

		if len(matched) == 0 {
			return
		} else {
			p, _ := strconv.Atoi(matched[1])
			if p == item.Priority {
				return true, err
			}
		}
	default:
		return true, err
	}
	return
}
