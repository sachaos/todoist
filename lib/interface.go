package todoist

type HaveID struct {
	ID int `json:"id"`
}

type HaveProjectID struct {
	ProjectID int `json:"project_id"`
}

type IDCarrier interface {
	GetID() int
}

type ContentCarrier interface {
	GetContent() string
}

type ProjectIDCarrier interface {
	GetProjectID() int
	GetProjectName(Projects) string
}

func (carrier HaveID) GetID() int {
	return carrier.ID
}

func (carrier HaveProjectID) GetProjectID() int {
	return carrier.ProjectID
}

func (carrier HaveProjectID) GetProjectName(projects Projects) string {
	project, err := projects.FindByID(carrier.GetProjectID())
	if err != nil {
		return ""
	}
	return project.Name
}
