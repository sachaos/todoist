package todoist

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

func (projects Projects) FindByID(id int) (Project, interface{}) {
	for _, project := range projects {
		if project.ID == id {
			return project, nil
		}
	}
	return Project{}, FindFailed
}
