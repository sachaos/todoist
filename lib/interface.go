package todoist

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type HaveID struct {
	ID int `json:"id"`
}
type HaveIDs []HaveID

type HaveProjectID struct {
	ProjectID int `json:"project_id"`
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

func SearchByID(repo Repository, id int) (data IDCarrier, err error) {
	index := sort.Search(repo.Len(), func(i int) bool {
		return repo.At(i).GetID() >= id
	})
	if index < repo.Len() && repo.At(index).GetID() == id {
		return repo.At(index), nil
	} else {
		return nil, errors.New("Find Failed")
	}
}

func SearchByIDPrefix(repo Repository, prefix string) (id int, err error) {
	index := sort.Search(repo.Len(), func(i int) bool {
		return strings.HasPrefix(strconv.Itoa(repo.At(i).GetID()), prefix)
	})
	if index < repo.Len() {
		if index < repo.Len() - 1 {
			if strings.HasPrefix(strconv.Itoa(repo.At(index + 1).GetID()), prefix) {
				// Ambiguous prefix, return converted input instead
				return strconv.Atoi(prefix)
			}
		}
		return repo.At(index).GetID(), nil
	}
	return strconv.Atoi(prefix)
}

type HaveParentID struct {
	ParentID interface{} `json:"parent_id"`
}

type ParentIDCarrier interface {
	GetParentID() (int, error)
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

func SearchParents(repo Repository, child ParentIDCarrier) (data []ParentIDCarrier, err error) {
	parentId, err := child.GetParentID()
	if err != nil {
		return []ParentIDCarrier{}, nil
	}
	childParent, err := SearchByID(repo, parentId)
	if err != nil {
		return []ParentIDCarrier{}, err
	}
	childParents, err := SearchParents(repo, childParent.(ParentIDCarrier))
	if err != nil {
		return []ParentIDCarrier{}, err
	}
	return append(childParents, childParent.(ParentIDCarrier)), nil
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

func (carrier HaveProjectID) GetProjectID() int {
	return carrier.ProjectID
}

func (carrier HaveProjectID) GetProjectName(projects Projects) string {
	project, err := SearchByID(projects, carrier.GetProjectID())
	if err != nil {
		return ""
	}
	return project.(Project).Name
}
