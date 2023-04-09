package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
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

func GenerateColorHash(ids []string, colorList []color.Attribute) map[string]color.Attribute {
	colorHash := map[string]color.Attribute{}
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
	return color.BlueString(carrier.GetID())
}

func ContentPrefix(store *todoist.Store, item *todoist.Item, depth int, c *cli.Context) (prefix string) {
	if c.GlobalBool("indent") {
		prefix = prefix + strings.Repeat("    ", depth)
	}
	if c.GlobalBool("namespace") {
		parents := todoist.SearchItemParents(store, item)
		for _, parent := range parents {
			prefix = prefix + parent.Content + ":"
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
	var p int
	switch priority {
	case 1:
		p = 4
		priorityColor.Add(color.FgBlue).Add(color.BgBlack)
	case 2:
		p = 3
		priorityColor.Add(color.FgHiYellow).Add(color.BgBlack)
	case 3:
		p = 2
		priorityColor.Add(color.FgHiRed).Add(color.BgBlack)
	case 4:
		p = 1
		priorityColor.Add(color.FgWhite).Add(color.BgRed)
	}
	return priorityColor.SprintFunc()(fmt.Sprintf("p%d", p))
}

func ProjectFormat(id string, store *todoist.Store, projectColorHash map[string]color.Attribute, c *cli.Context) string {
	var prefix string
	var namePrefix string
	project := store.FindProject(id)
	if project == nil {
		// Accept unknown project ID
		return color.New(color.FgCyan).SprintFunc()("Unknown")
	}

	projectName := project.Name
	if c.GlobalBool("project-namespace") {
		parentProjects := todoist.SearchProjectParents(store, project)
		for _, project := range parentProjects {
			namePrefix = namePrefix + project.Name + ":"
		}
	}
	return prefix + color.New(projectColorHash[project.GetID()]).SprintFunc()("#"+namePrefix+projectName)
}

func SectionFormat(id string, store *todoist.Store, c *cli.Context) string {
	prefix := ""
	sectionName := ""
	section := store.FindSection(id)
	if section != nil {
		prefix = "/"
		sectionName = section.Name
	}
	return prefix + sectionName
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
