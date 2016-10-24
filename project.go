package main

type Project struct {
	Collapsed    int         `json:"collapsed"`
	Color        int         `json:"color"`
	HasMoreNotes bool        `json:"has_more_notes"`
	ID           int         `json:"id"`
	InboxProject bool        `json:"inbox_project"`
	Indent       int         `json:"indent"`
	IsArchived   int         `json:"is_archived"`
	IsDeleted    int         `json:"is_deleted"`
	ItemOrder    int         `json:"item_order"`
	Name         string      `json:"name"`
	ParentID     interface{} `json:"parent_id"`
	Shared       bool        `json:"shared"`
}
