package main

import (
	"github.com/fatih/color"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
	"strings"
	"time"
)

func ColorList() []color.Attribute {
	return []color.Attribute{
		color.FgHiRed,
		color.FgHiGreen,
		color.FgHiYellow,
		color.FgHiBlue,
		color.FgHiMagenta,
		color.FgHiCyan,
	}
}

func GenerateColorHash(ids []int, colorList []color.Attribute) map[int]color.Attribute {
	colorHash := map[int]color.Attribute{}
	colorNum := 0
	for _, id := range ids {
		var colorAttribute color.Attribute
		value, ok := colorHash[id]
		if ok {
			colorAttribute = value
		} else {
			colorAttribute = colorList[colorNum]
			colorHash[id] = colorAttribute
			colorNum = colorNum + 1
			if colorNum == len(colorList) {
				colorNum = 0
			}
		}
	}
	return colorHash
}

func IdFormat(carrier todoist.IDCarrier) string {
	return color.BlueString(strconv.Itoa(carrier.GetID()))
}

func ContentPrefix(node *Node, c *cli.Context) (prefix string) {
	if c.GlobalBool("indent") {
		prefix = prefix + strings.Repeat("    ", node.Value.GetIndent())
	}
	if c.GlobalBool("namespace") {
		for _, pnode := range node.Parents() {
			prefix = prefix + pnode.Value.(todoist.ContentCarrier).GetContent() + ":"
		}
	}
	return
}

func ContentFormat(item todoist.ContentCarrier) string {
	if todoist.HasURL(item) {
		return color.New(color.Underline).SprintFunc()(todoist.GetContentTitle(item))
	}
	return todoist.GetContentTitle(item)
}

func PriorityFormat(priority int) string {
	priorityColor := color.New(color.Bold)
	switch priority {
	case 4:
		priorityColor.Add(color.FgWhite).Add(color.BgRed)
	case 3:
		priorityColor.Add(color.FgHiRed).Add(color.BgBlack)
	case 2:
		priorityColor.Add(color.FgHiYellow).Add(color.BgBlack)
	default:
		priorityColor.Add(color.FgBlue).Add(color.BgBlack)
	}
	return priorityColor.SprintFunc()("p" + strconv.Itoa(priority))
}

func ProjectFormat(id int, projects todoist.Projects, projectColorHash map[int]color.Attribute, c *cli.Context) string {
	var prefix string
	var namePrefix string
	project, err := todoist.SearchByID(projects, id)
	if err != nil {
		panic(err)
	}

	projectName := project.(todoist.Project).Name
	if c.GlobalBool("project-namespace") {
		parentProjects, err := todoist.SearchParents(projects, project.(todoist.Project))
		if err != nil {
			panic(err)
		}
		for _, project := range parentProjects {
			namePrefix = namePrefix + project.(todoist.Project).Name + ":"
		}
	}
	return prefix + color.New(projectColorHash[project.GetID()]).SprintFunc()("#"+namePrefix+projectName)
}

func dueDateString(dueDate time.Time, allDay bool) string {
	if (dueDate == time.Time{}) {
		return ""
	}
	dueDate = dueDate.Local()
	if !allDay {
		return dueDate.Format(ShortDateTimeFormat)
	}
	return dueDate.Format(ShortDateFormat)
}

func DueDateFormat(dueDate time.Time, allDay bool) string {
	dueDateString := dueDateString(dueDate, allDay)
	duration := time.Since(dueDate)
	dueDateColor := color.New(color.Bold)
	if duration > 0 {
		dueDateColor.Add(color.FgWhite).Add(color.BgRed)
	} else if duration > -12*time.Hour {
		dueDateColor.Add(color.FgHiRed).Add(color.BgBlack)
	} else if duration > -24*time.Hour {
		dueDateColor.Add(color.FgHiYellow).Add(color.BgBlack)
	} else {
		dueDateColor.Add(color.FgHiBlue).Add(color.BgBlack)
	}
	return dueDateColor.SprintFunc()(dueDateString)
}

func completedDateString(completedDate time.Time) string {
	if (completedDate == time.Time{}) {
		return ""
	}
	completedDate = completedDate.Local()
	return completedDate.Format(ShortDateTimeFormat)
}

func CompletedDateFormat(completedDate time.Time) string {
	return completedDateString(completedDate)
}
