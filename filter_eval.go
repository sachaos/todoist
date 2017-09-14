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
		e := e.(StringExpr)
		return EvalAsPriority(e, item), err
	default:
		return true, err
	}
	return
}

func EvalAsPriority(e StringExpr, item todoist.Item) (result bool) {
	matched := priorityRegex.FindStringSubmatch(e.literal)
	if len(matched) == 0 {
		return false
	} else {
		p, _ := strconv.Atoi(matched[1])
		if p == item.Priority {
			return true
		}
	}
	return false
}
