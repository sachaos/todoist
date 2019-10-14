package main

import (
	"regexp"
	"strconv"
	"time"

	"github.com/sachaos/todoist/lib"
)

var priorityRegex = regexp.MustCompile("^p([1-4])$")

func Eval(e Expression, item todoist.AbstractItem, projects todoist.Projects, labels todoist.Labels) (result bool, err error) {
	result = false
	switch e.(type) {
	case BoolInfixOpExpr:
		e := e.(BoolInfixOpExpr)
		lr, err := Eval(e.left, item, projects, labels)
		rr, err := Eval(e.right, item, projects, labels)
		if err != nil {
			return false, nil
		}
		switch e.operator {
		case '&':
			return lr && rr, nil
		case '|':
			return lr || rr, nil
		}
	case ProjectExpr:
		e := e.(ProjectExpr)
		return EvalProject(e, item.GetProjectID(), projects), err
	case LabelExpr:
		e := e.(LabelExpr)
		return EvalLabel(e, item.GetLabelIDs(), labels), err
	case StringExpr:
		switch item.(type) {
		case *todoist.Item:
			item := item.(*todoist.Item)
			e := e.(StringExpr)
			return EvalAsPriority(e, item), err
		default:
			return false, nil
		}
	case DateExpr:
		e := e.(DateExpr)
		dateTime, allDay := item.DateTime()
		return EvalDate(e, dateTime, allDay), err
	case NotOpExpr:
		e := e.(NotOpExpr)
		r, err := Eval(e.expr, item, projects, labels)
		if err != nil {
			return false, nil
		}
		return !r, nil
	default:
		return true, err
	}
	return
}

func EvalDate(e DateExpr, itemDate time.Time, allDay bool) (result bool) {
	if (itemDate == time.Time{}) {
		if e.operation == NO_DUE_DATE {
			return true
		}
		return false
	}
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

func EvalAsPriority(e StringExpr, item *todoist.Item) (result bool) {
	matched := priorityRegex.FindStringSubmatch(e.literal)
	if len(matched) == 0 {
		return false
	} else {
		p, _ := strconv.Atoi(matched[1])
		if p == priorityMapping[item.Priority] {
			return true
		}
	}
	return false
}

func EvalProject(e ProjectExpr, projectID int, projects todoist.Projects) bool {
	for _, id := range projects.GetIDsByName(e.name, e.isAll) {
		if id == projectID {
			return true
		}
	}
	return false
}

func EvalLabel(e LabelExpr, labelIDs []int, labels todoist.Labels) bool {
	if e.name == "" {
		if len(labelIDs) == 0 {
			return true
		} else {
			return false
		}
	}

	labelID := labels.GetIDByName(e.name)
	if labelID == 0 {
		return false
	}

	for _, id := range labelIDs {
		if id == labelID {
			return true
		}
	}

	return false
}
