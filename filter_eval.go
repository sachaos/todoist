package main

import (
	"regexp"
	"strconv"
	"time"

	"github.com/sachaos/todoist/lib"
)

var priorityRegex = regexp.MustCompile("^p([1-4])$")

// FIXME make eval more abstract for reusing
func ComplEval(e Expression, item todoist.CompletedItem) (result bool, err error) {
	result = false
	switch e.(type) {
	case BoolInfixOpExpr:
		e := e.(BoolInfixOpExpr)
		lr, err := ComplEval(e.left, item)
		rr, err := ComplEval(e.right, item)
		if err != nil {
			return false, nil
		}
		switch e.operator {
		case '&':
			return lr && rr, nil
		case '|':
			return lr || rr, nil
		}
	case DateExpr:
		e := e.(DateExpr)
		return EvalDate(e, item.DateTime()), err
	case NotOpExpr:
		e := e.(NotOpExpr)
		r, err := ComplEval(e.expr, item)
		if err != nil {
			return false, nil
		}
		return !r, nil
	default:
		return true, err
	}
	return
}

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
	case DateExpr:
		e := e.(DateExpr)
		return EvalDate(e, item.DateTime()), err
	case NotOpExpr:
		e := e.(NotOpExpr)
		r, err := Eval(e.expr, item)
		if err != nil {
			return false, nil
		}
		return !r, nil
	default:
		return true, err
	}
	return
}

func EvalDate(e DateExpr, itemDate time.Time) (result bool) {
	if (itemDate == time.Time{}) {
		if e.operation == NO_DUE_DATE {
			return true
		}
		return false
	}
	allDay := e.allDay
	dueDate := e.datetime
	switch e.operation {
	case DUE_ON:
		var startDate, endDate time.Time
		if allDay {
			startDate = dueDate
			endDate = dueDate.AddDate(0, 0, 1)
			if itemDate.Equal(startDate) || (itemDate.After(startDate) && itemDate.Before(endDate)) {
				return true
			}
		}
		return false
	case DUE_BEFORE:
		if itemDate.Before(dueDate) {
			return true
		}
		return false
	case DUE_AFTER:
		endDateTime := dueDate
		if allDay {
			endDateTime = dueDate.AddDate(0, 0, 1).Add(-time.Duration(time.Microsecond))
		}
		if itemDate.After(endDateTime) {
			return true
		}
		return false
	default:
		return false
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
