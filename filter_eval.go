package main

import (
	"regexp"
	"strconv"
	"time"

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
	case DueDateExpr:
		e := e.(DueDateExpr)
		return EvalDueDate(e, item), err
	default:
		return true, err
	}
	return
}

func EvalDueDate(e DueDateExpr, item todoist.Item) (result bool) {
	itemDueDate := item.DueDateTime()
	allDay := e.allDay
	dueDate := e.datetime
	switch e.operation {
	case DUE_ON:
		var startDate, endDate time.Time
		if allDay {
			startDate = dueDate
			endDate = dueDate.AddDate(0, 0, 1)
			if itemDueDate.Equal(startDate) || (itemDueDate.After(startDate) && itemDueDate.Before(endDate)) {
				return true
			}
		}
		return false
	case DUE_BEFORE:
		if itemDueDate.Before(dueDate) {
			return true
		}
		return false
	case DUE_AFTER:
		endDateTime := dueDate
		if allDay {
			endDateTime = dueDate.AddDate(0, 0, 1).Add(-time.Duration(time.Microsecond))
		}
		if itemDueDate.After(endDateTime) {
			return true
		}
		return false
	default:
		return true
	}
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
