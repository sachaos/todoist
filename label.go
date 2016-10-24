package main

type Label struct {
	Color     int    `json:"color"`
	ID        int    `json:"id"`
	IsDeleted int    `json:"is_deleted"`
	ItemOrder int    `json:"item_order"`
	Name      string `json:"name"`
}

func FindByID(lables []Label, id int) (Label, interface{}) {
	for _, label := range lables {
		if label.ID == id {
			return label, nil
		}
	}
	return Label{}, "NotFound"
}
