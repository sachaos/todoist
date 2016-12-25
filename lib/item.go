package lib

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	linkRegex = regexp.MustCompile(`\[(.*)\]\((.*)\)`)
)

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
	CompletedDate string      `json:"completed_date"`
	MetaData      interface{} `json:"meta_data"`
	TaskID        int         `json:"task_id"`
}

func (item CompletedItem) CompletedDateTime() time.Time {
	t, _ := time.Parse(DateFormat, item.CompletedDate)
	return t
}

type CompletedItems []CompletedItem

type Item struct {
	BaseItem
	AllDay         bool        `json:"all_day"`
	AssignedByUID  int         `json:"assigned_by_uid"`
	Checked        int         `json:"checked"`
	Collapsed      int         `json:"collapsed"`
	DateAdded      string      `json:"date_added"`
	DateLang       string      `json:"date_lang"`
	DateString     string      `json:"date_string"`
	DayOrder       int         `json:"day_order"`
	DueDateUtc     string      `json:"due_date_utc"`
	HasMoreNotes   bool        `json:"has_more_notes"`
	InHistory      int         `json:"in_history"`
	Indent         int         `json:"indent"`
	IsArchived     int         `json:"is_archived"`
	IsDeleted      int         `json:"is_deleted"`
	ItemOrder      int         `json:"item_order"`
	LabelIDs       []int       `json:"labels"`
	ParentID       interface{} `json:"parent_id"`
	Priority       int         `json:"priority"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	SyncID         interface{} `json:"sync_id"`
}

type Items []Item

func (item Item) DueDateTime() time.Time {
	t, _ := time.Parse(DateFormat, item.DueDateUtc)
	return t
}

func (items Items) FindByID(id int) (Item, error) {
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return Item{}, FindFailed
}

func GetContentTitle(item ContentCarrier) string {
	return linkRegex.ReplaceAllString(item.GetContent(), "$1")
}

func GetContentURL(item ContentCarrier) string {
	if HasURL(item) {
		return linkRegex.ReplaceAllString(item.GetContent(), "$2")
	}
	return ""
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

func (item Item) MoveParam(to_project Project) interface{} {
	param := map[string]interface{}{
		"project_items": map[string][]int{
			strconv.Itoa(item.ProjectID): []int{item.ID},
		},
		"to_project": to_project.ID,
	}
	return param
}

func (item Item) LabelsString(labels Labels) string {
	label_names := make([]string, 0)
	for _, label_id := range item.LabelIDs {
		label, err := labels.FindByID(label_id)
		if err != nil {
			return "Error"
		}
		label_names = append(label_names, "@"+label.Name)
	}
	return strings.Join(label_names, ",")
}

func AddItem(item Item, token string) error {
	commands := Commands{
		NewCommand("item_add", item.AddParam()),
	}
	_, err := SyncRequest(commands.UrlValues(token))
	return err
}

func UpdateItem(item Item, token string) error {
	commands := Commands{
		NewCommand("item_update", item.UpdateParam()),
	}
	_, err := SyncRequest(commands.UrlValues(token))
	return err
}

func CloseItem(ids []int, token string) error {
	var commands Commands
	for _, id := range ids {
		command := NewCommand("item_close", map[string]interface{}{"id": id})
		commands = append(commands, command)
	}
	_, err := SyncRequest(commands.UrlValues(token))
	return err
}

func DeleteItem(ids []int, token string) error {
	commands := Commands{
		NewCommand("item_delete", map[string]interface{}{"ids": ids}),
	}
	_, err := SyncRequest(commands.UrlValues(token))
	return err
}

func MoveItem(item Item, to_project Project, token string) error {
	commands := Commands{
		NewCommand("item_move", item.MoveParam(to_project)),
	}
	_, err := SyncRequest(commands.UrlValues(token))
	return err
}
