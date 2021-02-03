package todoist

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
