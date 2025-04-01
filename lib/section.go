package todoist

import "strings"

type Section struct {
	HaveID
	HaveProjectID
	Collapsed    bool   `json:"collapsed"`
	Name         string `json:"name"`
	IsArchived   bool   `json:"is_archived"`
	IsDeleted    bool   `json:"is_deleted"`
	SectionOrder int    `json:"section_order"`
}

type Sections []Section

func (a Sections) GetIDsByName(name string) []string {
	var ids []string
	name = strings.ToLower(name)
	for _, sec := range a {
		if strings.Contains(strings.ToLower(sec.Name), name) {
			ids = append(ids, sec.ID)
		}
	}
	return ids
}
