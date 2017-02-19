package todoist

type Order struct {
	Num  int         `json:"num"`
	ID   int         `json:"id"`
	Data interface{} `json:"-"`
}
type Orders []Order

func (a Orders) Len() int           { return len(a) }
func (a Orders) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Orders) Less(i, j int) bool { return a[i].Num < a[j].Num }

type ItemOrder struct {
	Order
	ProjectOrder int `json:"project_order"`
}

type ItemOrders []ItemOrder

func (a ItemOrders) Len() int      { return len(a) }
func (a ItemOrders) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ItemOrders) Less(i, j int) bool {
	if a[i].ProjectOrder == a[j].ProjectOrder {
		return a[i].Num < a[j].Num
	} else {
		return a[i].ProjectOrder < a[j].ProjectOrder
	}
}
