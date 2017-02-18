package todoist

import (
	"errors"
	"sort"
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

func (projects Projects) SearchParents(project *Project) (data []*Project, err error) {
	parentId := project.ParentID
	if parentId == nil {
		return []*Project{}, nil
	}
	projectParent, err := projects.SearchByID(int(parentId.(float64)))
	if err != nil {
		return []*Project{}, err
	}
	projectParents, err := projects.SearchParents(projectParent)
	if err != nil {
		return []*Project{}, err
	}
	return append(projectParents, projectParent), nil
}

func (projects Projects) SearchByID(id int) (data *Project, err error) {
	index := sort.Search(len(projects), func(i int) bool {
		return projects[i].ID >= id
	})
	if index < len(projects) && projects[index].ID == id {
		return &projects[index], nil
	} else {
		return nil, errors.New("Find Failed")
	}
}

func (projects Projects) FindByID(id int) (Project, interface{}) {
	for _, project := range projects {
		if project.ID == id {
			return project, nil
		}
	}
	return Project{}, FindFailed
}
