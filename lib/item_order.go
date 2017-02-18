package todoist

type ItemOrder struct {
	Num  int         `json:"num"`
	ID   int         `json:"id"`
	Data interface{} `json:"-"`
}

type ItemOrders []ItemOrder

func (a ItemOrders) Len() int           { return len(a) }
func (a ItemOrders) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ItemOrders) Less(i, j int) bool { return a[i].Num < a[j].Num }
