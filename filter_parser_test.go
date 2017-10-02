package main

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

// Test ...
func TestFilter(t *testing.T) {
	assert.Equal(t, nil, Filter(""), "the should be equal")
}

func TestPriorityFilter(t *testing.T) {
	assert.Equal(t, StringExpr{literal: "p1"}, Filter("p1"), "the should be equal")
}

func TestBoolInfixFilter(t *testing.T) {
	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '|',
			right:    StringExpr{literal: "p2"},
		},
		Filter("p1 | p2"), "the should be equal")

	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '&',
			right:    StringExpr{literal: "p2"},
		},
		Filter("p1 & p2"), "the should be equal")

	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '&',
			right: BoolInfixOpExpr{
				left:     StringExpr{literal: "p2"},
				operator: '|',
				right:    StringExpr{literal: "p3"},
			},
		},
		Filter("p1 & (p2 | p3 )"), "the should be equal")
}

func setNow(t time.Time) {
	now = func() time.Time { return t }
}

func TestDateTimeFilter(t *testing.T) {
	timeNow := time.Now()
	setNow(timeNow)

	assert.Equal(t,
		time.Date(2017, time.October, 5, 0, 0, 0, 0, time.Local),
		Filter("10/5/2017"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), time.January, 3, 0, 0, 0, 0, time.Local),
		Filter("Jan 3"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), time.August, 8, 0, 0, 0, 0, time.Local),
		Filter("8 August"), "the should be equal")

	assert.Equal(t,
		time.Date(2020, time.February, 10, 0, 0, 0, 0, time.Local),
		Filter("10 Feb 2020"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), time.May, 16, 0, 0, 0, 0, time.Local),
		Filter("16/05"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 16, 0, 0, 0, time.Local),
		Filter("16:00"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 16, 10, 3, 0, time.Local),
		Filter("16:10:03"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 15, 0, 0, 0, time.Local),
		Filter("3pm"), "the should be equal")

	assert.Equal(t,
		time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 7, 0, 0, 0, time.Local),
		Filter("7am"), "the should be equal")

	assert.Equal(t,
		time.Date(2020, time.February, 10, 15, 0, 0, 0, time.Local),
		Filter("10 Feb 2020 3pm"), "the should be equal")
}
