package views

type Item struct {
	mp *Map
	name string
	image string
	kind string
	pos Position
}

func NewItem(m *Map) *Item {
	return &Item{
		mp:m,
		name: "",
		image: "",
		kind: "",
		pos:Position{
			x:0,
			y:0},
	}
}

