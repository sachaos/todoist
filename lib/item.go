package todoist

import (
	"context"
	"regexp"
	"strings"
	"time"
)

var (
	linkRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
)

const (
	RFC3339Date                 = "2006-01-02"
	RFC3339DateTime             = "2006-01-02T15:04:05"
	RFC3339DateTimeWithTimeZone = "2006-01-02T15:04:05Z07:00"
)

type Due struct {
	Date        string `json:"date"`
	TimeZone    string `json:"timezone"`
	IsRecurring bool   `json:"is_recurring"`
	String      string `json:"string"`
	Lang        string `json:"lang"`
}

type BaseItem struct {
	HaveID
	HaveProjectID
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

func (bitem BaseItem) GetContent() string {
	return bitem.Content
}

type CompletedItem struct {
	BaseItem
	CompletedData string      `json:"completed_date"`
	MetaData      interface{} `json:"meta_data"`
	TaskID        int         `json:"task_id"`
}

func (item CompletedItem) DateTime() time.Time {
	t, _ := time.Parse(time.RFC3339, item.CompletedData)
	return t
}

func (item CompletedItem) GetProjectID() int {
	return item.ProjectID
}

func (item CompletedItem) GetLabelIDs() []int {
	return []int{}
}

type CompletedItems []CompletedItem

type Item struct {
	BaseItem
	HaveParentID
	HaveIndent
	HaveSectionID
	ChildItem      *Item       `json:"-"`
	BrotherItem    *Item       `json:"-"`
	AllDay         bool        `json:"all_day"`
	AssignedByUID  int         `json:"assigned_by_uid"`
	Checked        int         `json:"checked"`
	Collapsed      int         `json:"collapsed"`
	DateAdded      string      `json:"date_added"`
	DateLang       string      `json:"date_lang"`
	DateString     string      `json:"date_string"`
	DayOrder       int         `json:"day_order"`
	Due            *Due        `json:"due"`
	HasMoreNotes   bool        `json:"has_more_notes"`
	InHistory      int         `json:"in_history"`
	IsArchived     int         `json:"is_archived"`
	IsDeleted      int         `json:"is_deleted"`
	ItemOrder      int         `json:"item_order"`
	LabelIDs       []int       `json:"labels"`
	Priority       int         `json:"priority"`
	AutoReminder   bool        `json:"auto_reminder"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	SyncID         interface{} `json:"sync_id"`
}

type Items []Item

func (a Items) Len() int           { return len(a) }
func (a Items) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Items) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Items) At(i int) IDCarrier { return a[i] }

func (item Item) DateTime() time.Time {
	var date string

	if item.Due == nil {
		date = ""
	} else {
		date = item.Due.Date
	}
	//2020-03-03T14:00:00
	//2020-01-17T23:00:00Z
	t, err := time.ParseInLocation(RFC3339DateTimeWithTimeZone, date, time.Local)
	if err != nil {
		t, err = time.ParseInLocation(RFC3339DateTime, date, time.Local)
	}
	if err != nil {
		t, _ = time.ParseInLocation(RFC3339Date, date, time.Local)
	}
	return t
}

func (item Item) GetProjectID() int {
	return item.ProjectID
}

func (item Item) GetLabelIDs() []int {
	return item.LabelIDs
}

// interface for Eval actions
type AbstractItem interface {
	DateTime() time.Time
	GetProjectID() int
	GetLabelIDs() []int
}

func GetContentTitle(item ContentCarrier) string {
	return linkRegex.ReplaceAllString(item.GetContent(), "$1")
}

func GetContentURL(item ContentCarrier) []string {
	if HasURL(item) {
		matches := linkRegex.FindAllStringSubmatch(item.GetContent(), -1)
		if matches != nil {
			urls := make([]string, len(matches))
			for i, match := range matches {
				urls[i] = match[2]
			}
			return urls
		}
	}
	return []string{}
}

func HasURL(item ContentCarrier) bool {
	return linkRegex.MatchString(item.GetContent())
}

func (item Item) AddParam() interface{} {
	param := map[string]interface{}{}
	if item.Content != "" {
		param["content"] = item.Content
	}
	if item.DateString != "" {
		param["date_string"] = item.DateString
	}
	if len(item.LabelIDs) != 0 {
		param["labels"] = item.LabelIDs
	}
	if item.Priority != 0 {
		param["priority"] = item.Priority
	}
	if item.ProjectID != 0 {
		param["project_id"] = item.ProjectID
	}
	param["auto_reminder"] = item.AutoReminder

	return param
}

func (item Item) UpdateParam() interface{} {
	param := map[string]interface{}{}
	if item.ID != 0 {
		param["id"] = item.ID
	}
	if item.Content != "" {
		param["content"] = item.Content
	}
	if item.DateString != "" {
		param["date_string"] = item.DateString
	}
	// TODO: more cool
	if item.DateString == "null" {
		param["date_string"] = ""
	}
	if len(item.LabelIDs) != 0 {
		param["labels"] = item.LabelIDs
	}
	if item.Priority != 0 {
		param["priority"] = item.Priority
	}
	return param
}

func (item *Item) MoveParam(projectId int) interface{} {
	param := map[string]interface{}{
		"id":         item.ID,
		"project_id": projectId,
	}
	return param
}

func (item Item) LabelsString(store *Store) string {
	var b strings.Builder
	for i, labelId := range item.LabelIDs {
		label := store.FindLabel(labelId)
		b.WriteString("@" + label.Name)
		if i < len(item.LabelIDs)-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

func (c *Client) AddItem(ctx context.Context, item Item) error {
	commands := Commands{
		NewCommand("item_add", item.AddParam()),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) UpdateItem(ctx context.Context, item Item) error {
	commands := Commands{
		NewCommand("item_update", item.UpdateParam()),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) CloseItem(ctx context.Context, ids []int) error {
	var commands Commands
	for _, id := range ids {
		command := NewCommand("item_close", map[string]interface{}{"id": id})
		commands = append(commands, command)
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) DeleteItem(ctx context.Context, ids []int) error {
	var commands Commands
	for _, id := range ids {
		command := NewCommand("item_delete", map[string]interface{}{"id": id})
		commands = append(commands, command)
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) MoveItem(ctx context.Context, item *Item, projectId int) error {
	commands := Commands{
		NewCommand("item_move", item.MoveParam(projectId)),
	}
	return c.ExecCommands(ctx, commands)
}
