package todoist

import (
	"errors"
)

type HaveID struct {
	ID int `json:"id"`
}

type HaveParentID struct {
	ParentID interface{} `json:"parent_id"`
}

type HaveProjectID struct {
	ProjectID int `json:"project_id"`
}

type HaveIndent struct {
	Indent int `json:"indent"`
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

func (carrier HaveIndent) GetIndent() int {
	return carrier.Indent
}

func (carrier HaveParentID) GetParentID() (int, error) {
	switch yi := carrier.ParentID.(type) {
	case int:
		return yi, nil
	case float64:
		return int(yi), nil
	default:
		return 0, errors.New("Parent ID is null")
	}
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
