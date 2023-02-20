package todoist

type Label struct {
	HaveID
	Color     string `json:"color"`
	IsDeleted bool   `json:"is_deleted"`
	ItemOrder int    `json:"item_order"`
	Name      string `json:"name"`
}

type Labels []Label

func (a Labels) Len() int           { return len(a) }
func (a Labels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Labels) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Labels) At(i int) IDCarrier { return a[i] }

func (a Labels) GetIDByName(name string) string {
	for _, label := range a {
		if label.Name == name {
			return label.ID
		}
	}
	return ""
}
