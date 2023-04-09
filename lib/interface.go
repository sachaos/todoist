package todoist

import (
	"errors"
)

type HaveID struct {
	ID string `json:"id"`
}
type HaveIDs []HaveID

type HaveProjectID struct {
	ProjectID string `json:"project_id"`
}

type HaveSectionID struct {
	SectionID string `json:"section_id"`
}

type HaveIndent struct {
	Indent int `json:"indent"`
}

type IDCarrier interface {
	GetID() string
}
type Repository interface {
	Len() int
	At(int) IDCarrier
}

type HaveParentID struct {
	ParentID *string `json:"parent_id"`
}

type ParentIDCarrier interface {
	GetParentID() (string, error)
}

func (carrier HaveParentID) GetParentID() (string, error) {
	if carrier.ParentID == nil {
		return "", errors.New("Parent ID is null")
	}
	return *carrier.ParentID, nil
}

func SearchProjectParents(store *Store, project *Project) []*Project {
	if project.ParentID == nil {
		return []*Project{}
	}

	parentProject := store.FindProject(*project.ParentID)
	return append(SearchProjectParents(store, parentProject), parentProject)
}

func SearchItemParents(store *Store, item *Item) []*Item {
	if item.ParentID == nil {
		return []*Item{}
	}

	parentItem := store.FindItem(*item.ParentID)
	return append(SearchItemParents(store, parentItem), parentItem)
}

type ContentCarrier interface {
	GetContent() string
}

type ProjectIDCarrier interface {
	GetProjectID() string
}

func (carrier HaveID) GetID() string {
	return carrier.ID
}

func (carrier HaveIndent) GetIndent() int {
	return carrier.Indent
}

func (carrier HaveProjectID) GetProjectID() string {
	return carrier.ProjectID
}
