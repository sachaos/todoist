package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

var priorityMapping = map[int]int{
	1: 4,
	2: 3,
	3: 2,
	4: 1,
}

func Add(c *cli.Context) error {
	client := GetClient(c)

	item := todoist.Item{}
	if !c.Args().Present() {
		return CommandFailed
	}

	item.Content = c.Args().First()
	item.Priority = priorityMapping[c.Int("priority")]
	item.ProjectID = c.Int("project-id")
	if item.ProjectID == 0 {
		item.ProjectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}
	item.LabelIDs = func(str string) []int {
		stringIDs := strings.Split(str, ",")
		ids := []int{}
		for _, stringID := range stringIDs {
			id, err := strconv.Atoi(stringID)
			if err != nil {
				continue
			}
			ids = append(ids, id)
		}
		return ids
	}(c.String("label-ids"))

	item.DateString = c.String("date")
	item.AutoReminder = c.Bool("reminder")

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	return Sync(c)
}

var (
	tomorrowIdentHash = map[string]bool{
		"tomorrow": true,
		"tom":      true,
	}
)

func isDateIdent(token string) bool {
	if _, ok := MonthIdentHash[token]; ok {
		return true
	} else if _, ok := TodayIdentHash[token]; ok {
		return true
	} else if _, ok := tomorrowIdentHash[token]; ok {
		return true
	}
	return false
}

func QuickAdd(c *cli.Context) error {
	client := GetClient(c)
	item := todoist.Item{}

	content := c.Args().First()

	var s scanner.Scanner
	s.Init(strings.NewReader(content))
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() == "#" {
			s.Scan()
			project := s.TokenText()
			projectID := client.Store.Projects.GetIDByName(project)
			item.ProjectID = projectID
			content = strings.Replace(content, fmt.Sprintf("#%s", project), "", 1)
		} else if s.TokenText() == "@" {
			s.Scan()
			label := s.TokenText()
			labelID := client.Store.Labels.GetIDByName(label)
			if labelID == 0 {
				log.Printf("Label '%s' is not known", label)
			} else {
				item.LabelIDs = append(item.LabelIDs, labelID)
				content = strings.Replace(content, fmt.Sprintf("@%s", label), "", 1)
			}
		} else if _, ok := TodayIdentHash[s.TokenText()]; ok {
			item.DateString = "today"
			content = strings.Replace(content, s.TokenText(), "", 1)
		} else if _, ok := tomorrowIdentHash[s.TokenText()]; ok {
			item.DateString = "tomorrow"
			content = strings.Replace(content, s.TokenText(), "", 1)
		}
	}

	item.Content = content

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	if err := Sync(c); err != nil {
		return err
	}

	return nil
}
