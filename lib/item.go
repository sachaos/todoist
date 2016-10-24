package lib

import (
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

func LabelsString(item Item, labels []Label) string {
	label_names := make([]string, 0)
	for _, label_id := range item.LabelIDs {
		label, err := FindByID(labels, label_id)
		if err != nil {
			return "Error"
		}
		label_names = append(label_names, "@"+label.Name)
	}
	return strings.Join(label_names, ",")
}
