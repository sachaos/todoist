package todoist

import (
	"context"
	"strings"
)

type Project struct {
	HaveID
	HaveParentID
	HaveIndent
	Collapsed      bool     `json:"collapsed"`
	Color          string   `json:"color"`
	HasMoreNotes   bool     `json:"has_more_notes"`
	InboxProject   bool     `json:"inbox_project"`
	IsArchived     bool     `json:"is_archived"`
	IsDeleted      bool     `json:"is_deleted"`
	ItemOrder      int      `json:"item_order"`
	Name           string   `json:"name"`
	Shared         bool     `json:"shared"`
	ChildProject   *Project `json:"-"`
	BrotherProject *Project `json:"-"`
	ViewStyle      string   `json:"view_style"`
}

type Projects []Project

func (a Projects) Len() int           { return len(a) }
func (a Projects) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Projects) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Projects) At(i int) IDCarrier { return a[i] }

func (a Projects) GetIDByName(name string) string {
	for _, pjt := range a {
		if pjt.Name == name {
			return pjt.GetID()
		}
	}
	return ""
}

func (a Projects) GetIDsByName(name string, isAll bool) []string {
	ids := []string{}
	name = strings.ToLower(name)
	for _, pjt := range a {
		if strings.Contains(strings.ToLower(pjt.Name), name) {
			ids = append(ids, pjt.ID)
			if isAll {
				parentID := pjt.ID
				// Find all children which has the project as parent
				ids = append(ids, parentID)
				ids = append(ids, childProjectIDs(parentID, a)...)
			}
		}
	}
	return ids
}

func childProjectIDs(parentId string, projects Projects) []string {
	ids := []string{}
	for _, pjt := range projects {
		id, err := pjt.GetParentID()
		if err != nil {
			continue
		}

		if id == parentId {
			ids = append(ids, pjt.ID)
			ids = append(ids, childProjectIDs(pjt.ID, projects)...)
		}
	}
	return ids
}

func (project Project) AddParam() interface{} {
	param := map[string]interface{}{}
	if project.Name != "" {
		param["name"] = project.Name
	}
	if project.Color != "" {
		param["color"] = project.Color
	}
	//TODO: ParentID
	if project.ItemOrder != 0 {
		param["child_order"] = project.ItemOrder
	}
	//TODO: IsFavorite
	return param
}

func (c *Client) AddProject(ctx context.Context, project Project) error {
	commands := Commands{
		NewCommand("project_add", project.AddParam()),
	}
	return c.ExecCommands(ctx, commands)
}
