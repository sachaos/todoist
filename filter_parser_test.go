package main

import (
	"testing"

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

func TestDateTimeFilter(t *testing.T) {
	assert.Equal(t,
		SpecificDateTimeExpr{
			year:  2017,
			month: 10,
			day:   5,
		},
		Filter("10/5/2017"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			month: 1,
			day:   3,
		},
		Filter("Jan 3"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			month: 8,
			day:   8,
		},
		Filter("8 August"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			year:  2020,
			month: 2,
			day:   10,
		},
		Filter("10 Feb 2020"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			month: 5,
			day:   16,
		},
		Filter("16/05"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			hour:   16,
			minute: 00,
		},
		Filter("16:00"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			hour:   16,
			minute: 10,
			second: 3,
		},
		Filter("16:10:03"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			hour: 15,
		},
		Filter("3pm"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			hour: 7,
		},
		Filter("7am"), "the should be equal")

	assert.Equal(t,
		SpecificDateTimeExpr{
			year:  2020,
			month: 2,
			day:   10,
			hour:  15,
		},
		Filter("10 Feb 2020 3pm"), "the should be equal")
}
