package lib

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Item struct {
	AllDay         bool        `json:"all_day"`
	AssignedByUID  int         `json:"assigned_by_uid"`
	Checked        int         `json:"checked"`
	Collapsed      int         `json:"collapsed"`
	Content        string      `json:"content"`
	DateAdded      string      `json:"date_added"`
	DateLang       string      `json:"date_lang"`
	DateString     string      `json:"date_string"`
	DayOrder       int         `json:"day_order"`
	DueDateUtc     interface{} `json:"due_date_utc"`
	HasMoreNotes   bool        `json:"has_more_notes"`
	ID             int         `json:"id"`
	InHistory      int         `json:"in_history"`
	Indent         int         `json:"indent"`
	IsArchived     int         `json:"is_archived"`
	IsDeleted      int         `json:"is_deleted"`
	ItemOrder      int         `json:"item_order"`
	LabelIDs       []int       `json:"labels"`
	ParentID       interface{} `json:"parent_id"`
	Priority       int         `json:"priority"`
	ProjectID      int         `json:"project_id"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	SyncID         interface{} `json:"sync_id"`
	UserID         int         `json:"user_id"`
}

func (item Item) AddParam() interface{} {
	param := map[string]interface{}{}
	if item.Content != "" {
		param["content"] = item.Content
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
	if len(item.LabelIDs) != 0 {
		param["labels"] = item.LabelIDs
	}
	if item.Priority != 0 {
		param["priority"] = item.Priority
	}
	return param
}

func (item Item) ProjectString(projects Projects) string {
	project, err := projects.FindByID(item.ProjectID)
	if err != nil {
		return ""
	}
	return "#" + project.Name
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
	var commands []Command
	command := Command{
		UUID:   uuid.NewV4().String(),
		TempID: uuid.NewV4().String(),
		Type:   "item_add",
		Args:   item.AddParam(),
	}
	commands = append(commands, command)
	commands_text, err := json.Marshal(commands)
	if err != nil {
		return PostFailed
	}
	_, err = http.PostForm(
		"https://todoist.com/API/v7/sync",
		url.Values{"token": {token}, "commands": {string(commands_text)}},
	)
	if err != nil {
		return PostFailed
	}
	return nil
}

func UpdateItem(item Item, token string) error {
	var commands []Command
	command := Command{
		UUID:   uuid.NewV4().String(),
		TempID: uuid.NewV4().String(),
		Type:   "item_update",
		Args:   item.UpdateParam(),
	}
	commands = append(commands, command)
	commands_text, err := json.Marshal(commands)
	if err != nil {
		return PostFailed
	}
	_, err = http.PostForm(
		"https://todoist.com/API/v7/sync",
		url.Values{"token": {token}, "commands": {string(commands_text)}},
	)
	if err != nil {
		return PostFailed
	}
	return nil
}

func CloseItem(ids []int, token string) error {
	var commands []Command
	var command Command
	for _, id := range ids {
		command = Command{
			UUID:   uuid.NewV4().String(),
			TempID: uuid.NewV4().String(),
			Type:   "item_close",
		}
		command.Args = map[string]interface{}{"id": id}
		commands = append(commands, command)
	}
	commands_text, err := json.Marshal(commands)
	if err != nil {
		return PostFailed
	}
	_, err = http.PostForm(
		"https://todoist.com/API/v7/sync",
		url.Values{"token": {token}, "commands": {string(commands_text)}},
	)
	if err != nil {
		return PostFailed
	}
	return nil
}

func MoveItem(item Item, token string) error {
	var commands []Command
	var command Command
	command = Command{
		UUID:   uuid.NewV4().String(),
		TempID: uuid.NewV4().String(),
		Type:   "item_move",
	}
	command.Args = map[string]map[string][]int{
		"project_items": map[string][]int{
			strconv.Itoa(item.ProjectID): []int{item.ID},
		},
	}
	commands = append(commands, command)
	commands_text, err := json.Marshal(commands)
	if err != nil {
		return PostFailed
	}
	_, err = http.PostForm(
		"https://todoist.com/API/v7/sync",
		url.Values{"token": {token}, "commands": {string(commands_text)}},
	)
	if err != nil {
		return PostFailed
	}
	return nil
}
