package todoist

type Project struct {
	HaveID
	Collapsed    int         `json:"collapsed"`
	Color        int         `json:"color"`
	HasMoreNotes bool        `json:"has_more_notes"`
	InboxProject bool        `json:"inbox_project"`
	Indent       int         `json:"indent"`
	IsArchived   int         `json:"is_archived"`
	IsDeleted    int         `json:"is_deleted"`
	ItemOrder    int         `json:"item_order"`
	Name         string      `json:"name"`
	ParentID     interface{} `json:"parent_id"`
	Shared       bool        `json:"shared"`
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
