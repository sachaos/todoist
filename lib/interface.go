package todoist

import (
	"errors"
)

type HaveID struct {
	ID int `json:"id"`
}
type HaveIDs []HaveID

type HaveProjectID struct {
	ProjectID int `json:"project_id"`
}

// HaveSectionID defines the ID of the section that item belongs to, it can be emtpy(nul)
type HaveSectionID struct {
	SectionID interface{} `json:"section_id"`
}

type HaveIndent struct {
	Indent int `json:"indent"`
}

type IDCarrier interface {
	GetID() int
}
type Repository interface {
	Len() int
	At(int) IDCarrier
}

type HaveParentID struct {
	ParentID *int `json:"parent_id"`
}

type ParentIDCarrier interface {
	GetParentID() (int, error)
}

func (carrier HaveParentID) GetParentID() (int, error) {
	if carrier.ParentID == nil {
		return 0, errors.New("Parent ID is null")
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
	GetProjectID() int
}

func (carrier HaveID) GetID() int {
	return carrier.ID
}

func (carrier HaveIndent) GetIndent() int {
	return carrier.Indent
}

func (carrier HaveProjectID) GetProjectID() int {
	return carrier.ProjectID
}
