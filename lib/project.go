package todoist

import (
	"strings"
)

type Project struct {
	HaveID
	HaveParentID
	HaveIndent
	Collapsed    int    `json:"collapsed"`
	Color        int    `json:"color"`
	HasMoreNotes bool   `json:"has_more_notes"`
	InboxProject bool   `json:"inbox_project"`
	IsArchived   int    `json:"is_archived"`
	IsDeleted    int    `json:"is_deleted"`
	ItemOrder    int    `json:"item_order"`
	Name         string `json:"name"`
	Shared       bool   `json:"shared"`
}

type Projects []Project

func (a Projects) Len() int           { return len(a) }
func (a Projects) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Projects) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Projects) At(i int) IDCarrier { return a[i] }

func (a Projects) GetIDByName(name string) int {
	for _, pjt := range a {
		if pjt.Name == name {
			return pjt.GetID()
		}
	}
	return 0
}

func (a Projects) GetIDsByName(name string, isAll bool) []int {
	ids := []int{}
	for _, pjt := range a {
		if strings.Contains(pjt.Name, name) {
			ids = append(ids, pjt.ID)
			if isAll {
				parentID := pjt.ID
				// Find all children which has the project as parent
				ids = append(ids, parentID)
				for _, id := range childProjectIDs(parentID, a) {
					ids = append(ids, id)
				}
			}
		}
	}
	return ids
}

func childProjectIDs(parentId int, projects Projects) []int {
	ids := []int{}
	for _, pjt := range projects {
		id, err := pjt.GetParentID()
		if err != nil {
			continue
		}

		if id == parentId {
			ids = append(ids, pjt.ID)
			for _, id := range childProjectIDs(pjt.ID, projects) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
