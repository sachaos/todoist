package todoist

type Label struct {
	HaveID
	Color     int    `json:"color"`
	IsDeleted int    `json:"is_deleted"`
	ItemOrder int    `json:"item_order"`
	Name      string `json:"name"`
}

type Labels []Label

func (labels Labels) FindByID(id int) (Label, interface{}) {
	for _, label := range labels {
		if label.ID == id {
			return label, nil
		}
	}
	return Label{}, FindFailed
}
